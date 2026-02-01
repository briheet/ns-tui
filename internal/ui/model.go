package ui

import (
	"time"

	"github.com/briheet/ns-tui/internal/api"
	"github.com/briheet/ns-tui/internal/hm"
	"github.com/briheet/ns-tui/internal/models"
	"github.com/briheet/ns-tui/internal/styles"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// renderCache holds pre-rendered strings that only change on window resize.
type renderCache struct {
	banner      string
	bannerWidth int

	helpInsert string
	helpNormal string
	helpDetail string
	helpWidth  int
}

// Model represents the application state
type Model struct {
	textInput             textinput.Model
	packages              []models.Package
	cursor                int
	scrollOffset          int
	loading               bool
	err                   error
	width                 int
	height                int
	lastQuery             string
	searchTimer           *time.Timer
	mode                  models.Mode
	selectedPackage       *models.Package
	apiClient             *api.Client
	selectedInstallMethod int // 0-3 for the 4 install methods
	spinner               spinner.Model
	toastMessage          string
	toastVisible          bool
	showHelp              bool
	selectedTab           int       // 0=Nixpkgs, 1=Home Manager, 2=Pacman
	tabQueries            [3]string // Saved search text per tab
	cache                 renderCache
	// Home Manager state
	hmOptions         []models.HMOption // All loaded HM options
	hmSearchResults   []models.HMOption // Current HM search results
	hmLoaded          bool              // Whether HM options are loaded
	hmLoading         bool              // Whether HM fetch is in progress
	showHMPrompt      bool              // Show the fetch prompt modal
	hmPromptSelection int               // 0=Yes, 1=No
	hmCursor          int               // Cursor for HM results
	hmScrollOffset    int               // Scroll offset for HM results
	hmLastQuery       string            // Last HM search query
	hmErr             error             // HM-specific error
	// Home Manager detail state
	selectedHMOption      *models.HMOption       // Currently viewed HM option
	hmDetailHistory       []models.HmDetailEntry // Navigation stack for back-traversal
	hmRelatedOptions      []models.HMOption      // Sibling options for current selection
	hmRelatedCursor       int                    // Cursor in related options list
	hmRelatedScrollOffset int                    // Scroll offset for related options
}

// NewModel creates a new application model
func NewModel() Model {
	return Model{
		textInput:    models.NewTextInput(),
		packages:     []models.Package{},
		mode:         models.InsertMode,
		apiClient:    api.NewClient(),
		spinner:      models.NewSpinner(),
	}
}

// buildRenderCache pre-renders static UI sections that only change on resize.
func (m Model) buildRenderCache() renderCache {
	center := func(s string) string {
		return lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(s)
	}

	helpInsert := center(styles.HelpStyle.Render(
		"esc: normal mode • ↑/↓: navigate • tab: switch source • ?: help • q: quit"))
	helpNormal := center(styles.HelpStyle.Render(
		"i: insert mode • j/k: navigate • enter/space: details • g/G: top/bottom • tab: switch source • ?: help • q: quit"))
	helpDetail := center(styles.HelpStyle.Render(
		"j/k: navigate related • enter/space: view option • esc/b: back • ?: help • q: quit"))

	return renderCache{
		banner:      m.buildBannerHeader(),
		bannerWidth: m.width,
		helpInsert:  helpInsert,
		helpNormal:  helpNormal,
		helpDetail:  helpDetail,
		helpWidth:   m.width,
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

// checkHMCache checks if the HM options cache exists and loads it if so
func checkHMCache() tea.Cmd {
	return func() tea.Msg {
		if !hm.CacheExists() {
			return hmCacheCheckMsg{exists: false}
		}
		options, err := hm.LoadFromCache()
		return hmCacheCheckMsg{exists: true, options: options, err: err}
	}
}

// fetchHMOptions runs nix build and caches the result
func fetchHMOptions() tea.Cmd {
	return func() tea.Msg {
		options, err := hm.FetchAndCache()
		return hmFetchResultMsg{options: options, err: err}
	}
}

// searchHMOptions performs in-memory search over loaded HM options
func searchHMOptions(options []models.HMOption, query string) tea.Cmd {
	return func() tea.Msg {
		results := hm.Search(options, query, 50)
		return hmSearchResultMsg{results: results}
	}
}
