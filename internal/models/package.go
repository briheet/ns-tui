package models

// Package represents a NixOS package with all its metadata
type Package struct {
	Name            string
	Version         string
	Description     string
	AttrName        string
	AttrSet         string
	LongDescription string
	License         string
	LicenseSet      []string
	HomepageLinks   []string
	Platforms       []string
	Programs        []string
	Maintainers     []any
	MaintainersSet  []string
	Teams           []any
	TeamsSet        []string
	Outputs         []string
	MainProgram     string
	DefaultOutput   string
	Position        string
	System          string
	Hydra           map[string]any
}
