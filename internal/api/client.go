package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/briheet/ns-tui/internal/models"
)

const (
	elasticURL   = "https://search.nixos.org/backend"
	elasticIndex = "latest-44-nixos-unstable"
	username     = "aWVSALXpZv"
	password     = "X8gPHnzL52wFEekuxsfQ9cSh"
)

// Client handles communication with the NixOS search backend
type Client struct {
	httpClient *http.Client
	baseURL    string
	index      string
	auth       string
}

// NewClient creates a new API client
func NewClient() *Client {
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    elasticURL,
		index:      elasticIndex,
		auth:       auth,
	}
}

// esLicense represents a single license entry from the Elasticsearch response.
type esLicense struct {
	FullName string `json:"fullName"`
	URL      string `json:"url"`
}

// esSource is the typed representation of an Elasticsearch _source document.
type esSource struct {
	PName          string         `json:"package_pname"`
	PVersion       string         `json:"package_pversion"`
	Description    string         `json:"package_description"`
	AttrName       string         `json:"package_attr_name"`
	AttrSet        string         `json:"package_attr_set"`
	LongDesc       string         `json:"package_longDescription"`
	License        []esLicense    `json:"package_license"`
	LicenseSet     []string       `json:"package_license_set"`
	Homepage       []string       `json:"package_homepage"`
	Platforms      []string       `json:"package_platforms"`
	Programs       []string       `json:"package_programs"`
	Maintainers    []any          `json:"package_maintainers"`
	MaintainersSet []string       `json:"package_maintainers_set"`
	Teams          []any          `json:"package_teams"`
	TeamsSet       []string       `json:"package_teams_set"`
	Outputs        []string       `json:"package_outputs"`
	MainProgram    string         `json:"package_mainProgram"`
	DefaultOutput  string         `json:"package_default_output"`
	Position       string         `json:"package_position"`
	System         string         `json:"package_system"`
	Hydra          map[string]any `json:"package_hydra"`
}

// SearchPackages searches for packages matching the query
func (c *Client) SearchPackages(query string) ([]models.Package, error) {
	// Build Elasticsearch query
	esQuery := map[string]any{
		"from": 0,
		"size": 50,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"dis_max": map[string]any{
							"queries": []any{
								map[string]any{
									"multi_match": map[string]any{
										"query": query,
										"fields": []string{
											"package_attr_name^9",
											"package_pname^6",
											"package_attr_name.*^5.5",
											"package_pname.*^5.3",
											"package_programs^9",
											"package_programs.*^5.3",
										},
										"type":      "best_fields",
										"fuzziness": "AUTO",
									},
								},
								map[string]any{
									"wildcard": map[string]any{
										"package_attr_name": map[string]any{
											"value":            "*" + query + "*",
											"case_insensitive": true,
										},
									},
								},
							},
						},
					},
				},
				"should": []any{
					// Boost exact matches significantly
					map[string]any{
						"term": map[string]any{
							"package_attr_name": map[string]any{
								"value": query,
								"boost": 100,
							},
						},
					},
					map[string]any{
						"term": map[string]any{
							"package_pname": map[string]any{
								"value": query,
								"boost": 80,
							},
						},
					},
					// Boost prefix matches
					map[string]any{
						"prefix": map[string]any{
							"package_attr_name": map[string]any{
								"value": query,
								"boost": 50,
							},
						},
					},
				},
				"filter": []any{
					map[string]any{
						"term": map[string]any{
							"type": "package",
						},
					},
				},
			},
		},
		"sort": []any{
			"_score",
			map[string]any{
				"package_attr_name": "desc",
			},
			map[string]any{
				"package_pversion": "desc",
			},
		},
	}

	jsonData, err := json.Marshal(esQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/%s/_search", c.baseURL, c.index)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication
	req.Header.Add("Authorization", "Basic "+c.auth)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("elasticsearch returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response using streaming decoder
	var esResponse struct {
		Hits struct {
			Hits []struct {
				Source esSource `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&esResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract packages with deduplication
	packages := make([]models.Package, 0, len(esResponse.Hits.Hits))
	seen := make(map[string]bool)

	for _, hit := range esResponse.Hits.Hits {
		src := hit.Source

		// Create unique key for deduplication
		key := src.AttrName + ":" + src.PVersion
		if seen[key] {
			continue // Skip duplicates
		}
		seen[key] = true

		pkg := models.Package{
			Name:            src.PName,
			Version:         src.PVersion,
			Description:     src.Description,
			AttrName:        src.AttrName,
			AttrSet:         src.AttrSet,
			LongDescription: src.LongDesc,
			License:         formatLicenses(src.License),
			LicenseSet:      src.LicenseSet,
			HomepageLinks:   src.Homepage,
			Platforms:       src.Platforms,
			Programs:        src.Programs,
			Maintainers:     src.Maintainers,
			MaintainersSet:  src.MaintainersSet,
			Teams:           src.Teams,
			TeamsSet:        src.TeamsSet,
			Outputs:         src.Outputs,
			MainProgram:     src.MainProgram,
			DefaultOutput:   src.DefaultOutput,
			Position:        src.Position,
			System:          src.System,
			Hydra:           src.Hydra,
		}
		if pkg.Name == "" {
			pkg.Name = pkg.AttrName
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// formatLicenses joins license full names into a comma-separated string.
func formatLicenses(licenses []esLicense) string {
	names := make([]string, 0, len(licenses))
	for _, l := range licenses {
		if l.FullName != "" {
			names = append(names, l.FullName)
		}
	}
	return strings.Join(names, ", ")
}
