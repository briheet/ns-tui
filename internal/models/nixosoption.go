package models

// NixOSOption represents a single NixOS configuration option from the ES index
type NixOSOption struct {
	Name        string   // Dotted option path, e.g. "services.nginx.enable"
	Description string   // Human-readable description (may contain HTML)
	Type        string   // Nix type, e.g. "boolean", "null or (submodule)"
	Default     *string  // Default value text (nil if not set)
	Example     *string  // Example value text (nil if not set)
	Source      string   // Source file path, e.g. "nixos/modules/services/web-apps/mediawiki.nix"
	Loc         []string // Path components, e.g. ["services", "nginx", "enable"]
}

// NixOSDetailEntry stores one level of the NixOS option detail navigation stack
type NixOSDetailEntry struct {
	Option       NixOSOption
	Related      []NixOSOption
	Cursor       int
	ScrollOffset int
}
