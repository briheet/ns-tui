package hm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/briheet/ns-tui/internal/models"
)

const (
	cacheDirName  = "ns-tui"
	cacheFileName = "home-manager-options.json"
	nixFlakeRef   = "github:nix-community/home-manager#docs-json"
)

// CachePath returns the full path to the HM options cache file
func CachePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("could not determine cache directory: %w", err)
	}
	return filepath.Join(cacheDir, cacheDirName, cacheFileName), nil
}

// CacheExists checks whether the HM options JSON exists in the cache
func CacheExists() bool {
	path, err := CachePath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// FetchAndCache runs `nix build` to produce the docs-json output,
// copies the resulting JSON file to the cache directory, and returns parsed options.
func FetchAndCache() ([]models.HMOption, error) {
	// Run nix build to get the output path
	cmd := exec.Command("nix", "build", nixFlakeRef, "--no-link", "--print-out-paths")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nix build failed: %w", err)
	}
	nixOutPath := strings.TrimSpace(string(output))

	// The output contains options.json at share/doc/home-manager/options.json
	sourceFile := filepath.Join(nixOutPath, "share", "doc", "home-manager", "options.json")

	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read nix build output: %w", err)
	}

	// Ensure cache directory exists and write cache file
	cachePath, err := CachePath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create cache dir: %w", err)
	}
	if err := os.WriteFile(cachePath, data, 0o644); err != nil {
		return nil, fmt.Errorf("failed to write cache file: %w", err)
	}

	return ParseOptions(data)
}

// LoadFromCache reads and parses the cached HM options JSON
func LoadFromCache() ([]models.HMOption, error) {
	cachePath, err := CachePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}
	return ParseOptions(data)
}

// rawOption represents a single option entry in the JSON
type rawOption struct {
	Description  string            `json:"description"`
	Type         string            `json:"type"`
	Default      *json.RawMessage  `json:"default"`
	Example      *json.RawMessage  `json:"example"`
	Declarations []json.RawMessage `json:"declarations"`
	Loc          []string          `json:"loc"`
	ReadOnly     bool              `json:"readOnly"`
}

// ParseOptions parses the Home Manager options JSON into a slice of HMOption.
// The JSON is a flat dict: { "programs.git.enable": { ... }, ... }
func ParseOptions(data []byte) ([]models.HMOption, error) {
	var raw map[string]rawOption
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse HM options JSON: %w", err)
	}

	options := make([]models.HMOption, 0, len(raw))
	for name, entry := range raw {
		opt := models.HMOption{
			Name:        name,
			Description: entry.Description,
			Type:        entry.Type,
			Default:     parseHMValue(entry.Default),
			Example:     parseHMValue(entry.Example),
			Loc:         entry.Loc,
			ReadOnly:    entry.ReadOnly,
		}

		// Parse declarations
		for _, rawDecl := range entry.Declarations {
			// Try as object with name+url
			var decl models.HMDeclaration
			if err := json.Unmarshal(rawDecl, &decl); err == nil && decl.URL != "" {
				opt.Declarations = append(opt.Declarations, decl)
				continue
			}
			// Fallback: try as plain string
			var s string
			if err := json.Unmarshal(rawDecl, &s); err == nil {
				opt.Declarations = append(opt.Declarations, models.HMDeclaration{Name: s})
			}
		}

		options = append(options, opt)
	}

	// Sort by name for deterministic ordering
	sort.Slice(options, func(i, j int) bool {
		return options[i].Name < options[j].Name
	})

	return options, nil
}

// parseHMValue extracts a displayable text string from an HM option value.
// Handles both {_type: "literalExpression", text: "..."} objects and plain values.
func parseHMValue(raw *json.RawMessage) *string {
	if raw == nil {
		return nil
	}
	// Try as typed expression: {"_type": "literalExpression", "text": "..."}
	var typed struct {
		Type string `json:"_type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(*raw, &typed); err == nil && typed.Text != "" {
		return &typed.Text
	}
	// Fallback: try as plain string
	var s string
	if err := json.Unmarshal(*raw, &s); err == nil {
		return &s
	}
	// Last resort: use raw JSON
	text := string(*raw)
	return &text
}

// Search performs case-insensitive substring matching on option names and descriptions.
// Returns up to limit results, scored: exact > prefix > name-contains > description-only.
func Search(options []models.HMOption, query string, limit int) []models.HMOption {
	if query == "" {
		return nil
	}

	lowerQuery := strings.ToLower(query)

	type scored struct {
		opt   models.HMOption
		score int
	}
	var results []scored

	for _, opt := range options {
		lowerName := strings.ToLower(opt.Name)
		lowerDesc := strings.ToLower(opt.Description)

		nameMatch := strings.Contains(lowerName, lowerQuery)
		descMatch := strings.Contains(lowerDesc, lowerQuery)

		if !nameMatch && !descMatch {
			continue
		}

		score := 100 // description-only match
		if nameMatch {
			score = 10 // name contains
			if strings.HasPrefix(lowerName, lowerQuery) {
				score = 1 // prefix match
			}
			if lowerName == lowerQuery {
				score = 0 // exact match
			}
		}
		results = append(results, scored{opt: opt, score: score})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].score != results[j].score {
			return results[i].score < results[j].score
		}
		return results[i].opt.Name < results[j].opt.Name
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	out := make([]models.HMOption, len(results))
	for i, r := range results {
		out[i] = r.opt
	}
	return out
}
