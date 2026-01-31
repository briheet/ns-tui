package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var esResponse struct {
		Hits struct {
			Hits []struct {
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(body, &esResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract packages with deduplication
	packages := make([]models.Package, 0, len(esResponse.Hits.Hits))
	seen := make(map[string]bool)

	for _, hit := range esResponse.Hits.Hits {
		src := hit.Source
		attrName := getString(src, "package_attr_name")
		version := getString(src, "package_pversion")

		// Create unique key for deduplication
		key := attrName + ":" + version
		if seen[key] {
			continue // Skip duplicates
		}
		seen[key] = true

		pkg := models.Package{
			Name:            getString(src, "package_pname"),
			Version:         version,
			Description:     getString(src, "package_description"),
			AttrName:        attrName,
			AttrSet:         getString(src, "package_attr_set"),
			LongDescription: getString(src, "package_longDescription"),
			License:         getString(src, "package_license"),
			LicenseSet:      getArray(src, "package_license_set"),
			HomepageLinks:   getArray(src, "package_homepage"),
			Platforms:       getArray(src, "package_platforms"),
			Programs:        getArray(src, "package_programs"),
			Maintainers:     getArray(src, "package_maintainers"),
			MaintainersSet:  getArray(src, "package_maintainers_set"),
			Teams:           getArray(src, "package_teams"),
			TeamsSet:        getArray(src, "package_teams_set"),
			Outputs:         getArray(src, "package_outputs"),
			MainProgram:     getString(src, "package_mainProgram"),
			DefaultOutput:   getString(src, "package_default_output"),
			Position:        getString(src, "package_position"),
			System:          getString(src, "package_system"),
			Hydra:           getMap(src, "package_hydra"),
		}
		if pkg.Name == "" {
			pkg.Name = pkg.AttrName
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

func getString(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getArray(m map[string]any, key string) []any {
	if val, ok := m[key]; ok {
		if arr, ok := val.([]any); ok {
			return arr
		}
	}
	return []any{}
}

func getMap(m map[string]any, key string) map[string]any {
	if val, ok := m[key]; ok {
		if mapVal, ok := val.(map[string]any); ok {
			return mapVal
		}
	}
	return nil
}
