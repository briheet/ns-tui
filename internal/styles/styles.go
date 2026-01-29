package styles

import "github.com/charmbracelet/lipgloss"

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
