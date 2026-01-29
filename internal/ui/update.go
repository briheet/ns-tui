package ui

import (
	"fmt"
	"time"

	"github.com/briheet/ns-tui/internal/models"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle help overlay
		if msg.String() == "?" {
			m.showHelp = !m.showHelp
			return m, nil
		}

		// Close help with esc
		if m.showHelp && msg.String() == "esc" {
			m.showHelp = false
			return m, nil
		}

		// Close tab message with esc or enter
		if m.showTabMessage && (msg.String() == "esc" || msg.String() == "enter") {
			m.showTabMessage = false
			return m, nil
		}

		// Handle HM fetch prompt
		if m.showHMPrompt {
			switch msg.String() {
			case "j", "k", "down", "up", "tab", "shift+tab":
				m.hmPromptSelection = (m.hmPromptSelection + 1) % 2
				return m, nil
			case "enter":
				m.showHMPrompt = false
				if m.hmPromptSelection == 0 {
					// User chose Yes — start fetch
					m.hmLoading = true
					return m, fetchHMOptions()
				}
				// User chose No — go back to Nixpkgs
				m.selectedTab = 0
				m.textInput.Placeholder = "Search NixOS packages..."
				return m, nil
			case "esc":
				m.showHMPrompt = false
				m.selectedTab = 0
				m.textInput.Placeholder = "Search NixOS packages..."
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			}
			return m, nil
		}

		// Don't process other keys when help or tab message is shown
		if m.showHelp || m.showTabMessage {
			return m, nil
		}

		// Handle detail mode separately
		if m.mode == models.DetailMode {
			return m.handleDetailModeKeys(msg)
		}

		// Handle tab cycling in search view (not in detail mode)
		if msg.String() == "tab" {
			m.selectedTab = (m.selectedTab + 1) % 3
			return m, m.handleTabSwitch()
		}
		if msg.String() == "shift+tab" {
			m.selectedTab = (m.selectedTab - 1 + 3) % 3
			return m, m.handleTabSwitch()
		}

		// Handle keys based on mode
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.mode == models.InsertMode {
				m.mode = models.NormalMode
				m.textInput.Blur()
				return m, nil
			}
			// In normal mode, esc does nothing (use q to quit)
			return m, nil
		case "i":
			if m.mode == models.NormalMode {
				m.mode = models.InsertMode
				m.textInput.Focus()
				return m, textinput.Blink
			}
		case "j":
			if m.mode == models.NormalMode {
				if m.selectedTab == 0 && len(m.packages) > 0 {
					if m.cursor < len(m.packages)-1 {
						m.cursor++
					}
				} else if m.selectedTab == 1 && len(m.hmSearchResults) > 0 {
					if m.hmCursor < len(m.hmSearchResults)-1 {
						m.hmCursor++
					}
				}
				return m, nil
			}
		case "k":
			if m.mode == models.NormalMode {
				if m.selectedTab == 0 && len(m.packages) > 0 {
					if m.cursor > 0 {
						m.cursor--
					}
				} else if m.selectedTab == 1 && len(m.hmSearchResults) > 0 {
					if m.hmCursor > 0 {
						m.hmCursor--
					}
				}
				return m, nil
			}
		case "g":
			if m.mode == models.NormalMode {
				if m.selectedTab == 0 && len(m.packages) > 0 {
					m.cursor = 0
					m.scrollOffset = 0
				} else if m.selectedTab == 1 && len(m.hmSearchResults) > 0 {
					m.hmCursor = 0
					m.hmScrollOffset = 0
				}
				return m, nil
			}
		case "G":
			if m.mode == models.NormalMode {
				if m.selectedTab == 0 && len(m.packages) > 0 {
					m.cursor = len(m.packages) - 1
				} else if m.selectedTab == 1 && len(m.hmSearchResults) > 0 {
					m.hmCursor = len(m.hmSearchResults) - 1
				}
				return m, nil
			}
		case "enter", " ":
			if m.mode == models.InsertMode {
				// Switch from Insert to Normal mode (only on Enter, not space)
				if msg.String() == "enter" {
					m.mode = models.NormalMode
					m.textInput.Blur()
					return m, nil
				}
			}
			if m.mode == models.NormalMode {
				if m.selectedTab == 0 && len(m.packages) > 0 {
					m.mode = models.DetailMode
					m.selectedPackage = &m.packages[m.cursor]
					return m, nil
				}
				// HM detail mode not implemented yet — no-op
			}
		case "down":
			if m.mode == models.InsertMode {
				if m.selectedTab == 0 && len(m.packages) > 0 {
					if m.cursor < len(m.packages)-1 {
						m.cursor++
					}
				} else if m.selectedTab == 1 && len(m.hmSearchResults) > 0 {
					if m.hmCursor < len(m.hmSearchResults)-1 {
						m.hmCursor++
					}
				}
				return m, nil
			}
		case "up":
			if m.mode == models.InsertMode {
				if m.selectedTab == 0 && len(m.packages) > 0 {
					if m.cursor > 0 {
						m.cursor--
					}
				} else if m.selectedTab == 1 && len(m.hmSearchResults) > 0 {
					if m.hmCursor > 0 {
						m.hmCursor--
					}
				}
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case searchResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.packages = msg.packages
			m.cursor = 0
			m.scrollOffset = 0
		}
		return m, nil

	case hmCacheCheckMsg:
		if msg.exists && msg.err == nil {
			m.hmOptions = msg.options
			m.hmLoaded = true
			m.textInput.Placeholder = "Search Home Manager options..."
		} else if msg.exists && msg.err != nil {
			m.hmErr = msg.err
		} else {
			// Cache doesn't exist — show the fetch prompt
			m.showHMPrompt = true
			m.hmPromptSelection = 0
		}
		return m, nil

	case hmFetchResultMsg:
		m.hmLoading = false
		if msg.err != nil {
			m.hmErr = msg.err
		} else {
			m.hmOptions = msg.options
			m.hmLoaded = true
			m.textInput.Placeholder = "Search Home Manager options..."
		}
		return m, nil

	case hmSearchResultMsg:
		m.hmSearchResults = msg.results
		m.hmCursor = 0
		m.hmScrollOffset = 0
		m.loading = false
		return m, nil

	case clipboardMsg:
		// Show toast notification for copy success/failure
		if msg.success {
			m.toastMessage = fmt.Sprintf("✓ Copied: %s", msg.command)
			m.toastVisible = true
			// Hide toast after 2 seconds
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return hideToastMsg{}
			})
		} else {
			m.toastMessage = fmt.Sprintf("✗ Copy failed: %v", msg.err)
			m.toastVisible = true
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return hideToastMsg{}
			})
		}

	case hideToastMsg:
		m.toastVisible = false
		m.toastMessage = ""
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Update text input only in insert mode
	if m.mode == models.InsertMode {
		m.textInput, cmd = m.textInput.Update(msg)

		// Trigger search if input changed
		if m.textInput.Value() != m.lastQuery {
			m.lastQuery = m.textInput.Value()
			if m.lastQuery != "" {
				m.loading = true
				if m.selectedTab == 0 {
					// Nixpkgs: API search with debounce
					return m, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
						return performSearch(m.apiClient, m.lastQuery)()
					})
				} else if m.selectedTab == 1 && m.hmLoaded {
					// Home Manager: in-memory search with shorter debounce
					m.hmLastQuery = m.lastQuery
					return m, tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
						return searchHMOptions(m.hmOptions, m.hmLastQuery)()
					})
				}
			} else {
				if m.selectedTab == 0 {
					m.packages = []models.Package{}
					m.cursor = 0
					m.scrollOffset = 0
				} else if m.selectedTab == 1 {
					m.hmSearchResults = []models.HMOption{}
					m.hmCursor = 0
					m.hmScrollOffset = 0
				}
			}
		}
	}

	return m, cmd
}

// handleTabSwitch handles logic when switching tabs
func (m *Model) handleTabSwitch() tea.Cmd {
	switch m.selectedTab {
	case 0:
		// Nixpkgs tab
		m.textInput.Placeholder = "Search NixOS packages..."
		return nil
	case 1:
		// Home Manager tab
		if m.hmLoaded {
			m.textInput.Placeholder = "Search Home Manager options..."
			// If there's a query, re-search with HM data
			if m.textInput.Value() != "" {
				m.hmLastQuery = m.textInput.Value()
				return searchHMOptions(m.hmOptions, m.hmLastQuery)
			}
			return nil
		} else if m.hmLoading {
			// Already fetching, just show spinner
			return nil
		}
		// Not loaded — check cache
		return checkHMCache()
	default:
		// Pacman — still under development
		m.showTabMessage = true
		return nil
	}
}

// copyInstallCommand copies the selected installation command to clipboard
func (m Model) copyInstallCommand() tea.Cmd {
	if m.selectedPackage == nil {
		return nil
	}

	var command string
	switch m.selectedInstallMethod {
	case 0: // nix-shell
		command = fmt.Sprintf("nix-shell -p %s", m.selectedPackage.AttrName)
	case 1: // NixOS Config
		command = fmt.Sprintf("pkgs.%s", m.selectedPackage.AttrName)
	case 2: // nix-env
		command = fmt.Sprintf("nix-env -iA nixpkgs.%s", m.selectedPackage.AttrName)
	case 3: // nix profile
		command = fmt.Sprintf("nix profile install nixpkgs#%s", m.selectedPackage.AttrName)
	}

	return func() tea.Msg {
		err := clipboard.WriteAll(command)
		if err != nil {
			return clipboardMsg{success: false, err: err}
		}
		return clipboardMsg{success: true, command: command}
	}
}

// handleDetailModeKeys handles key presses in detail mode
func (m Model) handleDetailModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "backspace", "b":
		m.mode = models.NormalMode
		m.selectedPackage = nil
		m.selectedInstallMethod = 0 // Reset selection
		return m, nil
	case "tab", "j":
		// Cycle through install methods (0-3)
		m.selectedInstallMethod = (m.selectedInstallMethod + 1) % 4
		return m, nil
	case "shift+tab", "k":
		// Cycle backwards
		m.selectedInstallMethod = (m.selectedInstallMethod - 1 + 4) % 4
		return m, nil
	case "enter", " ":
		// Copy the selected installation command to clipboard
		return m, m.copyInstallCommand()
	}
	return m, nil
}
