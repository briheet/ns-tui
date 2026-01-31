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
	LicenseSet      []any
	HomepageLinks   []any
	Platforms       []any
	Programs        []any
	Maintainers     []any
	MaintainersSet  []any
	Teams           []any
	TeamsSet        []any
	Outputs         []any
	MainProgram     string
	DefaultOutput   string
	Position        string
	System          string
	Hydra           map[string]any
}
