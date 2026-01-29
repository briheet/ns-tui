package ui

import (
	"fmt"
	"time"

	"ns-tui/internal/models"

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

		// Don't process other keys when help is shown
		if m.showHelp {
			return m, nil
		}

		// Handle detail mode separately
		if m.mode == models.DetailMode {
			return m.handleDetailModeKeys(msg)
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
			return m, tea.Quit
		case "i":
			if m.mode == models.NormalMode {
				m.mode = models.InsertMode
				m.textInput.Focus()
				return m, textinput.Blink
			}
		case "j":
			if m.mode == models.NormalMode && len(m.packages) > 0 {
				if m.cursor < len(m.packages)-1 {
					m.cursor++
				}
				return m, nil
			}
		case "k":
			if m.mode == models.NormalMode && len(m.packages) > 0 {
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			}
		case "g":
			if m.mode == models.NormalMode && len(m.packages) > 0 {
				m.cursor = 0
				m.scrollOffset = 0
				return m, nil
			}
		case "G":
			if m.mode == models.NormalMode && len(m.packages) > 0 {
				m.cursor = len(m.packages) - 1
				return m, nil
			}
		case "enter":
			if m.mode == models.InsertMode {
				// Switch from Insert to Normal mode
				m.mode = models.NormalMode
				m.textInput.Blur()
				return m, nil
			}
			if m.mode == models.NormalMode && len(m.packages) > 0 {
				m.mode = models.DetailMode
				m.selectedPackage = &m.packages[m.cursor]
				return m, nil
			}
		case "down":
			if m.mode == models.InsertMode && len(m.packages) > 0 {
				if m.cursor < len(m.packages)-1 {
					m.cursor++
				}
				return m, nil
			}
		case "up":
			if m.mode == models.InsertMode && len(m.packages) > 0 {
				if m.cursor > 0 {
					m.cursor--
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
				// Debounce the search
				return m, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
					return performSearch(m.apiClient, m.lastQuery)()
				})
			} else {
				m.packages = []models.Package{}
				m.cursor = 0
				m.scrollOffset = 0
			}
		}
	}

	return m, cmd
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
		command = fmt.Sprintf("environment.systemPackages = [ pkgs.%s ];", m.selectedPackage.AttrName)
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
