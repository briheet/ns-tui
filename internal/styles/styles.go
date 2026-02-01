package styles

import (
	"strings"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

// Color palette - Catppuccin inspired
var (
	ColorPink        = lipgloss.Color("#FF6BCB")
	ColorBlue        = lipgloss.Color("#89B4FA")
	ColorPurple      = lipgloss.Color("#CBA6F7")
	ColorPinkLight   = lipgloss.Color("#F5C2E7")
	ColorCyan        = lipgloss.Color("#89DCEB")
	ColorYellow      = lipgloss.Color("#F9E2AF")
	ColorWhite       = lipgloss.Color("#CDD6F4")
	ColorGreen       = lipgloss.Color("#A6E3A1")
	ColorRed         = lipgloss.Color("#F38BA8")
	ColorTeal        = lipgloss.Color("#94E2D5")
	ColorGray        = lipgloss.Color("#6C7086")
	ColorDarkGray    = lipgloss.Color("#45475A")
	ColorBg          = lipgloss.Color("#1E1E2E")
	ColorBgHighlight = lipgloss.Color("#313244")
)

// Title and header styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPink).
			Background(ColorBg).
			Padding(0, 2).
			Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Italic(true).
			Align(lipgloss.Center).
			MarginBottom(1)
)

// Search box styles
var (
	SearchBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPinkLight).
		Padding(0, 1).
		Width(70)
)

// Mode indicator styles
var (
	InsertModeStyle = lipgloss.NewStyle().
			Bold(true).
			Background(ColorGreen).
			Foreground(ColorBg).
			Padding(0, 1)

	NormalModeStyle = lipgloss.NewStyle().
			Bold(true).
			Background(ColorBlue).
			Foreground(ColorBg).
			Padding(0, 1)
)

// Result item styles
var (
	ResultItemStyle = lipgloss.NewStyle().
			Padding(0, 1).
			MarginTop(1)

	SelectedItemStyle = lipgloss.NewStyle().
				Background(ColorBgHighlight).
				Foreground(ColorPurple).
				Padding(0, 1).
				MarginTop(1).
				Bold(true)

	PackageNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorCyan)

	VersionStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(ColorWhite)
)

// Status and message styles
var (
	LoadingStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true).
			Align(lipgloss.Center)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	CountStyle = lipgloss.NewStyle().
			Foreground(ColorTeal).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Align(lipgloss.Center).
			MarginTop(1)
)

// Detail view styles
var (
	// Note: DetailBoxStyle width is now set dynamically in renderDetailView
	DetailBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPinkLight).
			Padding(1, 2)

	DetailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPinkLight)

	DetailValueStyle = lipgloss.NewStyle().
				Foreground(ColorWhite)
)

// Separator style
var SeparatorStyle = lipgloss.NewStyle().
	Foreground(ColorDarkGray)

// Banner styles
var (
	BannerStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	BannerGrayStyle = lipgloss.NewStyle().
				Foreground(ColorGray)
)

// Tab styles
var (
	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(ColorPurple).
			Background(ColorBgHighlight).
			Bold(true).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(ColorGray).
				Padding(0, 2)

	TabSeparatorStyle = lipgloss.NewStyle().
				Foreground(ColorGray)
)

// Inline UI element styles
var (
	PositionStyle = lipgloss.NewStyle().
			Foreground(ColorGray)

	ToastStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Background(ColorDarkGray).
			Padding(0, 2).
			Bold(true)

	ScrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(ColorGray)

	NoResultsStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	HintStyle = lipgloss.NewStyle().
			Foreground(ColorTeal).
			Italic(true)

	HintGrayStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Italic(true)

	ProgramsStyle = lipgloss.NewStyle().
			Foreground(ColorTeal).
			Faint(true)

	WhiteTextStyle = lipgloss.NewStyle().
			Foreground(ColorWhite)
)

// Detail view inline styles
var (
	URLStyle = lipgloss.NewStyle().
			Foreground(ColorBlue)

	BreadcrumbSegmentStyle = lipgloss.NewStyle().
				Foreground(ColorCyan).
				Bold(true)

	BreadcrumbSepStyle = lipgloss.NewStyle().
				Foreground(ColorGray)

	ReadOnlyStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	PlatformSupportedStyle = lipgloss.NewStyle().
				Foreground(ColorGreen)

	PlatformUnsupportedStyle = lipgloss.NewStyle().
					Foreground(ColorRed)

	IndicatorStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	SelectedMethodStyle = lipgloss.NewStyle().
				Foreground(ColorPurple).
				Bold(true)

	NormalMethodStyle = lipgloss.NewStyle().
				Foreground(ColorWhite)

	CmdStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	CopyHelpStyle = lipgloss.NewStyle().
			Foreground(ColorTeal).
			Bold(true)

	CopiedHelpStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Bold(true)

	LeftBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorGray).
			Padding(0, 1)

	RightBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGreen).
			Padding(0, 1)
)

// Overlay styles
var (
	OverlayTitleStyle = lipgloss.NewStyle().
				Foreground(ColorPurple).
				Bold(true).
				Align(lipgloss.Center)

	OverlayMessageStyle = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Align(lipgloss.Center)

	OverlayInfoStyle = lipgloss.NewStyle().
				Foreground(ColorGray).
				Italic(true).
				Align(lipgloss.Center)

	OverlayBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple).
			Padding(2, 4).
			Width(60)
)

// Pacman placeholder styles
var (
	PacmanTitleStyle = lipgloss.NewStyle().
				Foreground(ColorYellow).
				Bold(true)

	PacmanMsgStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Italic(true)
)

// SetTheme reassigns all color and style variables from a Catppuccin flavor.
// Valid names: mocha, latte, frappe, macchiato. Defaults to mocha.
func SetTheme(name string) {
	name = strings.ToLower(name)
	var flavor catppuccin.Flavor
	switch name {
	case "latte":
		flavor = catppuccin.Latte
	case "frappe":
		flavor = catppuccin.Frappe
	case "macchiato":
		flavor = catppuccin.Macchiato
	default:
		flavor = catppuccin.Mocha
	}

	// Reassign color variables
	ColorPink = lipgloss.Color(flavor.Rosewater().Hex)
	ColorBlue = lipgloss.Color(flavor.Blue().Hex)
	ColorPurple = lipgloss.Color(flavor.Mauve().Hex)
	ColorPinkLight = lipgloss.Color(flavor.Pink().Hex)
	ColorCyan = lipgloss.Color(flavor.Sky().Hex)
	ColorYellow = lipgloss.Color(flavor.Yellow().Hex)
	ColorWhite = lipgloss.Color(flavor.Text().Hex)
	ColorGreen = lipgloss.Color(flavor.Green().Hex)
	ColorRed = lipgloss.Color(flavor.Red().Hex)
	ColorTeal = lipgloss.Color(flavor.Teal().Hex)
	ColorGray = lipgloss.Color(flavor.Overlay0().Hex)
	ColorDarkGray = lipgloss.Color(flavor.Surface1().Hex)
	ColorBg = lipgloss.Color(flavor.Base().Hex)
	ColorBgHighlight = lipgloss.Color(flavor.Surface0().Hex)

	// Rebuild all styles using updated colors
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPink).
		Background(ColorBg).
		Padding(0, 2).
		Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
		Foreground(ColorBlue).
		Italic(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	SearchBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPinkLight).
		Padding(0, 1).
		Width(70)

	InsertModeStyle = lipgloss.NewStyle().
		Bold(true).
		Background(ColorGreen).
		Foreground(ColorBg).
		Padding(0, 1)

	NormalModeStyle = lipgloss.NewStyle().
		Bold(true).
		Background(ColorBlue).
		Foreground(ColorBg).
		Padding(0, 1)

	ResultItemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		MarginTop(1)

	SelectedItemStyle = lipgloss.NewStyle().
		Background(ColorBgHighlight).
		Foreground(ColorPurple).
		Padding(0, 1).
		MarginTop(1).
		Bold(true)

	PackageNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorCyan)

	VersionStyle = lipgloss.NewStyle().
		Foreground(ColorYellow)

	DescriptionStyle = lipgloss.NewStyle().
		Foreground(ColorWhite)

	LoadingStyle = lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true).
		Align(lipgloss.Center)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ColorRed).
		Bold(true)

	CountStyle = lipgloss.NewStyle().
		Foreground(ColorTeal).
		Bold(true)

	HelpStyle = lipgloss.NewStyle().
		Foreground(ColorGray).
		Align(lipgloss.Center).
		MarginTop(1)

	DetailBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPinkLight).
		Padding(1, 2)

	DetailLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPinkLight)

	DetailValueStyle = lipgloss.NewStyle().
		Foreground(ColorWhite)

	SeparatorStyle = lipgloss.NewStyle().
		Foreground(ColorDarkGray)

	// Banner styles
	BannerStyle = lipgloss.NewStyle().
		Foreground(ColorBlue).
		Bold(true)

	BannerGrayStyle = lipgloss.NewStyle().
		Foreground(ColorGray)

	// Tab styles
	ActiveTabStyle = lipgloss.NewStyle().
		Foreground(ColorPurple).
		Background(ColorBgHighlight).
		Bold(true).
		Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
		Foreground(ColorGray).
		Padding(0, 2)

	TabSeparatorStyle = lipgloss.NewStyle().
		Foreground(ColorGray)

	// Inline UI element styles
	PositionStyle = lipgloss.NewStyle().
		Foreground(ColorGray)

	ToastStyle = lipgloss.NewStyle().
		Foreground(ColorGreen).
		Background(ColorDarkGray).
		Padding(0, 2).
		Bold(true)

	ScrollIndicatorStyle = lipgloss.NewStyle().
		Foreground(ColorGray)

	NoResultsStyle = lipgloss.NewStyle().
		Foreground(ColorYellow)

	HintStyle = lipgloss.NewStyle().
		Foreground(ColorTeal).
		Italic(true)

	HintGrayStyle = lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true)

	ProgramsStyle = lipgloss.NewStyle().
		Foreground(ColorTeal).
		Faint(true)

	WhiteTextStyle = lipgloss.NewStyle().
		Foreground(ColorWhite)

	// Detail view inline styles
	URLStyle = lipgloss.NewStyle().
		Foreground(ColorBlue)

	BreadcrumbSegmentStyle = lipgloss.NewStyle().
		Foreground(ColorCyan).
		Bold(true)

	BreadcrumbSepStyle = lipgloss.NewStyle().
		Foreground(ColorGray)

	ReadOnlyStyle = lipgloss.NewStyle().
		Foreground(ColorYellow)

	PlatformSupportedStyle = lipgloss.NewStyle().
		Foreground(ColorGreen)

	PlatformUnsupportedStyle = lipgloss.NewStyle().
		Foreground(ColorRed)

	IndicatorStyle = lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true)

	SelectedMethodStyle = lipgloss.NewStyle().
		Foreground(ColorPurple).
		Bold(true)

	NormalMethodStyle = lipgloss.NewStyle().
		Foreground(ColorWhite)

	CmdStyle = lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true)

	CopyHelpStyle = lipgloss.NewStyle().
		Foreground(ColorTeal).
		Bold(true)

	CopiedHelpStyle = lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true)

	LeftBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorGray).
		Padding(0, 1)

	RightBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorGreen).
		Padding(0, 1)

	// Overlay styles
	OverlayTitleStyle = lipgloss.NewStyle().
		Foreground(ColorPurple).
		Bold(true).
		Align(lipgloss.Center)

	OverlayMessageStyle = lipgloss.NewStyle().
		Foreground(ColorWhite).
		Align(lipgloss.Center)

	OverlayInfoStyle = lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true).
		Align(lipgloss.Center)

	OverlayBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPurple).
		Padding(2, 4).
		Width(60)

	// Pacman placeholder styles
	PacmanTitleStyle = lipgloss.NewStyle().
		Foreground(ColorYellow).
		Bold(true)

	PacmanMsgStyle = lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true)
}
