package ui

import (
	"time"

	"ns-tui/internal/api"
	"ns-tui/internal/models"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	textInput            textinput.Model
	packages             []models.Package
	cursor               int
	scrollOffset         int
	loading              bool
	err                  error
	width                int
	height               int
	lastQuery            string
	searchTimer          *time.Timer
	mode                 models.Mode
	selectedPackage      *models.Package
	apiClient            *api.Client
	selectedInstallMethod int // 0-3 for the 4 install methods
	spinner              spinner.Model
	toastMessage         string
	toastVisible         bool
	showHelp             bool
	selectedTab          int  // 0=Nixpkgs, 1=Home Manager, 2=Pacman
	showTabMessage       bool // Show "under development" message
}

// NewModel creates a new application model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Search NixOS packages..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 66

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		textInput: ti,
		packages:  []models.Package{},
		cursor:    0,
		loading:   false,
		mode:      models.InsertMode,
		apiClient: api.NewClient(),
		spinner:   s,
		toastVisible: false,
		showHelp:  false,
		selectedTab: 0,
		showTabMessage: false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

// performSearch executes the search in a goroutine
func performSearch(client *api.Client, query string) tea.Cmd {
	return func() tea.Msg {
		packages, err := client.SearchPackages(query)
		return searchResultMsg{packages: packages, err: err}
	}
}
