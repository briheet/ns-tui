package models

// HMOption represents a single Home Manager configuration option
type HMOption struct {
	Name         string          // Dotted option path, e.g. "programs.git.enable"
	Description  string          // Human-readable description
	Type         string          // Nix type, e.g. "boolean", "null or string"
	Default      *string         // Default value text (nil if not set)
	Example      *string         // Example value text (nil if not set)
	Declarations []HMDeclaration // Source file declarations
	Loc          []string        // Path components, e.g. ["programs", "git", "enable"]
	ReadOnly     bool            // Whether the option is read-only
}

// HMDeclaration represents a source file where an option is declared
type HMDeclaration struct {
	Name string // e.g. "<home-manager/modules/programs/git.nix>"
	URL  string // e.g. "https://github.com/nix-community/home-manager/blob/master/modules/programs/git.nix"
}
