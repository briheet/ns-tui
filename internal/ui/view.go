package ui

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/briheet/ns-tui/internal/models"
	"github.com/briheet/ns-tui/internal/styles"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m Model) View() string {
	// Show help overlay
	if m.showHelp {
		return m.renderHelpOverlay()
	}

	// Show tab message overlay
	if m.showTabMessage {
		return m.renderTabMessageOverlay()
	}

	// Show detailed package view
	if m.mode == models.DetailMode && m.selectedPackage != nil {
		return m.renderDetailView()
	}

	return m.renderSearchView()
}

// renderSearchView renders the main search view
func (m Model) renderSearchView() string {
	// Calculate heights
	footerHeight := 1 // help bar
	if m.toastVisible {
		footerHeight = 2 // toast + help bar
	}

	// HEADER SECTION
	var header strings.Builder
	header.WriteString("\n\n\n")

	// ASCII Banner
	bannerStyle := lipgloss.NewStyle().Foreground(styles.ColorBlue).Bold(true)
	bannerGrayStyle := lipgloss.NewStyle().Foreground(styles.ColorGray)

	bannerLines := []string{
		"'##::: ##::'######:::::::::::'########:'##::::'##:'####:",
		" ###:: ##:'##... ##::::::::::... ##..:: ##:::: ##:. ##::",
		" ####: ##: ##:::..:::::::::::::: ##:::: ##:::: ##:: ##::",
		" ## ## ##:. ######::'#######:::: ##:::: ##:::: ##:: ##::",
		" ##. ####::..... ##:........:::: ##:::: ##:::: ##:: ##::",
		" ##:. ###:'##::: ##::::::::::::: ##:::: ##:::: ##:: ##::",
		" ##::. ##:. ######:::::::::::::: ##::::. #######::'####:",
	}

	for _, line := range bannerLines {
		styledLine := bannerStyle.Render(line)
		header.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(styledLine))
		header.WriteString("\n")
	}

	dotsLine := bannerGrayStyle.Render("..::::..:::......:::::::::::::::..::::::.......:::.....::")
	header.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(dotsLine))
	header.WriteString("\n\n")

	subtitle := styles.SubtitleStyle.Render("Real-time package discovery with fuzzy search")
	header.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(subtitle))
	header.WriteString("\n\n")

	// Tabs
	activeTabStyle := lipgloss.NewStyle().
		Foreground(styles.ColorPurple).
		Background(lipgloss.Color("235")).
		Bold(true).
		Padding(0, 2)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(styles.ColorGray).
		Padding(0, 2)

	tabNames := []string{"Nixpkgs", "Home Manager", "Pacman"}
	var tabParts []string

	for i, name := range tabNames {
		if i > 0 {
			tabParts = append(tabParts, lipgloss.NewStyle().Foreground(styles.ColorGray).Render("│"))
		}
		if i == m.selectedTab {
			tabParts = append(tabParts, activeTabStyle.Render(name))
		} else {
			tabParts = append(tabParts, inactiveTabStyle.Render(name))
		}
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabParts...)
	header.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(tabs))
	header.WriteString("\n\n")

	// Search box with mode indicator
	var modeIndicator string
	if m.mode == models.InsertMode {
		modeIndicator = styles.InsertModeStyle.Render("-- INSERT --")
	} else {
		modeIndicator = styles.NormalModeStyle.Render("-- NORMAL --")
	}

	if len(m.packages) > 0 {
		position := lipgloss.NewStyle().
			Foreground(styles.ColorGray).
			Render(fmt.Sprintf(" [%d/%d]", m.cursor+1, len(m.packages)))
		modeIndicator = modeIndicator + position
	}

	searchBox := styles.SearchBoxStyle.Render(m.textInput.View())
	header.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(searchBox))
	header.WriteString("\n")
	header.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(modeIndicator))
	header.WriteString("\n")

	separator := styles.SeparatorStyle.Render(strings.Repeat("─", min(m.width, 80)))
	header.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(separator))
	header.WriteString("\n\n")

	headerContent := header.String()

	// RESULTS SECTION
	resultsContent := m.renderResults()

	// FOOTER SECTION (fixed at bottom)
	var footer strings.Builder

	// Toast (if visible)
	if m.toastVisible {
		toastStyle := lipgloss.NewStyle().
			Foreground(styles.ColorGreen).
			Background(lipgloss.Color("236")).
			Padding(0, 2).
			Bold(true)
		toast := toastStyle.Render(m.toastMessage)
		footer.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(toast))
		footer.WriteString("\n")
	}

	// Help bar
	var helpText string
	if m.mode == models.InsertMode {
		helpText = "esc: normal mode • ↑/↓: navigate • tab: switch source • ?: help • q: quit"
	} else {
		helpText = "i: insert mode • j/k: navigate • enter/space: details • g/G: top/bottom • tab: switch source • ?: help • q: quit"
	}
	help := styles.HelpStyle.Render(helpText)
	footer.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(help))

	footerContent := footer.String()

	// COMBINE ALL SECTIONS
	// Use lipgloss.Place to position footer at bottom
	mainAndResults := headerContent + resultsContent

	// Calculate how much space we need to fill to push footer to bottom
	// Count lines in mainAndResults
	mainLines := strings.Count(mainAndResults, "\n")
	footerLines := footerHeight

	// Fill remaining space
	availableHeight := m.height - footerLines - 1
	if mainLines < availableHeight {
		padding := availableHeight - mainLines
		for i := 0; i < padding; i++ {
			mainAndResults += "\n"
		}
	}

	return mainAndResults + "\n" + footerContent
}

// renderResults renders the search results
func (m Model) renderResults() string {
	var content strings.Builder

	// Loading indicator with spinner
	if m.loading {
		loading := fmt.Sprintf("%s Searching...", m.spinner.View())
		loadingStyled := styles.LoadingStyle.Render(loading)
		content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(loadingStyled))
		content.WriteString("\n")
		return content.String()
	}

	// Error message
	if m.err != nil {
		errorMsg := styles.ErrorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err))
		content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(errorMsg))
		content.WriteString("\n")
		return content.String()
	}

	// Results
	if len(m.packages) > 0 {
		// Calculate visible window first
		maxVisible := 8
		if m.height > 20 {
			maxVisible = 12
		}

		visibleCount := min(maxVisible, len(m.packages))
		count := styles.CountStyle.Render(fmt.Sprintf("📦 %d packages (showing %d)", len(m.packages), visibleCount))
		content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(count))
		content.WriteString("\n\n")

		// Ensure cursor is in view
		if m.cursor < m.scrollOffset {
			m.scrollOffset = m.cursor
		}
		if m.cursor >= m.scrollOffset+maxVisible {
			m.scrollOffset = m.cursor - maxVisible + 1
		}

		start := m.scrollOffset
		end := min(m.scrollOffset+maxVisible, len(m.packages))

		// Render packages
		for i := start; i < end; i++ {
			content.WriteString(m.renderPackageItem(i))
		}

		// Scroll indicators
		if m.scrollOffset > 0 {
			scrollUp := lipgloss.NewStyle().Foreground(styles.ColorGray).Render("⬆ More above")
			content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(scrollUp))
			content.WriteString("\n")
		}
		if end < len(m.packages) {
			scrollDown := lipgloss.NewStyle().Foreground(styles.ColorGray).Render(fmt.Sprintf("⬇ %d more below", len(m.packages)-end))
			content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(scrollDown))
			content.WriteString("\n")
		}
	} else if !m.loading && m.lastQuery != "" {
		noResults := lipgloss.NewStyle().
			Foreground(styles.ColorYellow).
			Render("No packages found. Try a different search term.")
		content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(noResults))
		content.WriteString("\n")
	} else if m.lastQuery == "" && !m.loading {
		hint := lipgloss.NewStyle().
			Foreground(styles.ColorTeal).
			Italic(true).
			Render("Type to search for NixOS packages...")
		content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(hint))
		content.WriteString("\n")
	}

	return content.String()
}

// renderPackageItem renders a single package item
func (m Model) renderPackageItem(index int) string {
	pkg := m.packages[index]

	cursor := "  "
	if m.cursor == index {
		cursor = "▶ "
	}

	name := styles.PackageNameStyle.Render(pkg.Name)
	version := styles.VersionStyle.Render(fmt.Sprintf("v%s", pkg.Version))
	desc := pkg.Description
	maxDescLen := 70
	if m.width > 100 {
		maxDescLen = 100
	}
	if len(desc) > maxDescLen {
		desc = desc[:maxDescLen-3] + "..."
	}
	desc = styles.DescriptionStyle.Render(desc)

	// Build line with name, version, and description
	var line string
	if len(pkg.Programs) > 0 {
		// Format programs list
		programsStyle := lipgloss.NewStyle().Foreground(styles.ColorTeal).Faint(true)
		var programsList []string
		maxPrograms := 5 // Show max 5 programs
		for i, prog := range pkg.Programs {
			if i >= maxPrograms {
				programsList = append(programsList, fmt.Sprintf("+%d more", len(pkg.Programs)-maxPrograms))
				break
			}
			if progStr, ok := prog.(string); ok {
				programsList = append(programsList, progStr)
			}
		}
		programsText := programsStyle.Render("📦 " + strings.Join(programsList, ", "))
		line = fmt.Sprintf("%s%s %s\n     %s\n     %s", cursor, name, version, desc, programsText)
	} else {
		line = fmt.Sprintf("%s%s %s\n     %s", cursor, name, version, desc)
	}

	var renderedLine string
	if m.cursor == index {
		renderedLine = styles.SelectedItemStyle.Render(line)
	} else {
		renderedLine = styles.ResultItemStyle.Render(line)
	}

	return lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(renderedLine) + "\n"
}

// renderDetailView renders the detailed package view
func (m Model) renderDetailView() string {
	var content strings.Builder

	pkg := m.selectedPackage

	// Calculate responsive box width (90% of screen, max 160)
	boxWidth := int(float64(m.width) * 0.9)
	if boxWidth > 160 {
		boxWidth = 160
	}

	// Center style for top 4 items (account for padding: 2 left + 2 right = 4)
	centerStyle := lipgloss.NewStyle().Width(boxWidth - 4).Align(lipgloss.Center)

	// Title
	title := styles.TitleStyle.Render(fmt.Sprintf("📦 %s", pkg.Name))
	content.WriteString(centerStyle.Render(title))
	content.WriteString("\n\n")

	// Version
	versionLine := styles.DetailLabelStyle.Render("Version: ") + styles.DetailValueStyle.Render(pkg.Version)
	content.WriteString(centerStyle.Render(versionLine))
	content.WriteString("\n\n")

	// Attribute Name
	attrLine := styles.DetailLabelStyle.Render("Attribute: ") + styles.DetailValueStyle.Render(pkg.AttrName)
	content.WriteString(centerStyle.Render(attrLine))
	content.WriteString("\n\n")

	// Description
	if pkg.Description != "" {
		descLabel := styles.DetailLabelStyle.Render("Description:")
		content.WriteString(centerStyle.Render(descLabel))
		content.WriteString("\n")
		descValue := styles.DetailValueStyle.Render(pkg.Description)
		content.WriteString(centerStyle.Render(descValue))
		content.WriteString("\n\n")
	}

	// Long Description (wrapped to box width)
	if pkg.LongDescription != "" {
		longDesc := stripHTML(pkg.LongDescription)
		content.WriteString(styles.DetailLabelStyle.Render("Long Description:\n"))
		// Wrap text to fit within box (account for padding and margins)
		wrappedDesc := wrapText(longDesc, boxWidth-8)
		content.WriteString(styles.DetailValueStyle.Render(wrappedDesc))
		content.WriteString("\n\n")
	}

	// License
	if pkg.License != "" {
		content.WriteString(styles.DetailLabelStyle.Render("License: "))
		content.WriteString(styles.DetailValueStyle.Render(pkg.License))
		content.WriteString("\n\n")
	}

	// Homepage
	if pkg.Homepage != "" {
		content.WriteString(styles.DetailLabelStyle.Render("Homepage: "))
		content.WriteString(styles.DetailValueStyle.Render(pkg.Homepage))
		content.WriteString("\n\n")
	}

	// Main Program
	if pkg.MainProgram != "" {
		content.WriteString(styles.DetailLabelStyle.Render("Main Program: "))
		content.WriteString(styles.DetailValueStyle.Render(pkg.MainProgram))
		content.WriteString("\n\n")
	}

	// Programs (show all with wrapping)
	if len(pkg.Programs) > 0 {
		content.WriteString(styles.DetailLabelStyle.Render("Programs: "))
		programsList := formatStringArray(pkg.Programs)
		// Wrap programs list to fit within box
		wrappedPrograms := wrapText(programsList, boxWidth-8)
		content.WriteString(styles.DetailValueStyle.Render(wrappedPrograms))
		content.WriteString("\n\n")
	}

	// Platform Compatibility Check
	if len(pkg.Platforms) > 0 {
		currentPlatform := getCurrentPlatform()
		isSupported := isPlatformSupported(pkg.Platforms, currentPlatform)

		content.WriteString(styles.DetailLabelStyle.Render("Platform: "))
		if isSupported {
			content.WriteString(lipgloss.NewStyle().Foreground(styles.ColorGreen).Render("✓ Supported on your platform"))
			content.WriteString(styles.DetailValueStyle.Render(fmt.Sprintf(" (%s)", currentPlatform)))
		} else {
			content.WriteString(lipgloss.NewStyle().Foreground(styles.ColorRed).Render("✗ Not supported on your platform"))
			content.WriteString(styles.DetailValueStyle.Render(fmt.Sprintf(" (%s)", currentPlatform)))
		}
		content.WriteString("\n\n")
	}

	// Installation Methods Header
	content.WriteString(styles.DetailLabelStyle.Render(fmt.Sprintf("How to install %s?\n\n", pkg.Name)))

	// Define methods and commands - ordered by ease of use
	methodNames := []string{"nix-shell", "NixOS Config", "nix-env", "nix profile"}
	commands := []string{
		fmt.Sprintf("nix-shell -p %s", pkg.AttrName),
		fmt.Sprintf("environment.systemPackages = [ pkgs.%s ];", pkg.AttrName),
		fmt.Sprintf("nix-env -iA nixpkgs.%s", pkg.AttrName),
		fmt.Sprintf("nix profile install nixpkgs#%s", pkg.AttrName),
	}

	// Build left menu - pad all items to same width for consistency
	leftWidth := 20  // Reduced for better balance
	rightWidth := 65 // Adjusted for better balance

	// Styles for menu items
	indicatorStyle := lipgloss.NewStyle().Foreground(styles.ColorGreen).Bold(true)
	selectedMethodStyle := lipgloss.NewStyle().Foreground(styles.ColorPurple).Bold(true)
	normalMethodStyle := lipgloss.NewStyle().Foreground(styles.ColorWhite)

	var menuItems []string
	for i, name := range methodNames {
		var line string
		if i == m.selectedInstallMethod {
			// Selected item: colored indicator + colored method name
			indicator := indicatorStyle.Render("→ ")
			methodName := selectedMethodStyle.Render(name)
			line = indicator + methodName
			// Calculate padding (accounting for ANSI codes)
			visibleLen := 2 + len(name) // "→ " + name
			if visibleLen < leftWidth {
				line = line + strings.Repeat(" ", leftWidth-visibleLen)
			}
		} else {
			// Normal item: spaces + normal method name
			methodName := normalMethodStyle.Render(name)
			line = "  " + methodName
			visibleLen := 2 + len(name)
			if visibleLen < leftWidth {
				line = line + strings.Repeat(" ", leftWidth-visibleLen)
			}
		}
		menuItems = append(menuItems, line)
	}
	leftContent := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Left box
	leftBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.ColorGray).
		Padding(0, 1)

	leftBox := leftBoxStyle.Render(leftContent)

	// Right box content - styled command
	selectedCmd := commands[m.selectedInstallMethod]

	// Style the command with green color and center it
	cmdStyle := lipgloss.NewStyle().Foreground(styles.ColorGreen).Bold(true)
	styledCmd := cmdStyle.Render(selectedCmd)

	// Center the command
	visibleCmdLen := len(selectedCmd)
	if visibleCmdLen < rightWidth {
		cmdLeftPad := (rightWidth - visibleCmdLen) / 2
		cmdRightPad := rightWidth - visibleCmdLen - cmdLeftPad
		styledCmd = strings.Repeat(" ", cmdLeftPad) + styledCmd + strings.Repeat(" ", cmdRightPad)
	} else if visibleCmdLen > rightWidth {
		// Truncate the original, then style
		selectedCmd = selectedCmd[:rightWidth]
		styledCmd = cmdStyle.Render(selectedCmd)
	}

	// Style help text - bold and visible
	helpStyle := lipgloss.NewStyle().Foreground(styles.ColorTeal).Bold(true)
	helpText := "Press Enter or Spacebar to copy"

	// Change text when copied successfully
	if m.toastVisible {
		helpStyle = lipgloss.NewStyle().Foreground(styles.ColorGreen).Bold(true)
		helpText = "✓ Copied successfully"
	}

	styledHelp := helpStyle.Render(helpText)

	helpVisibleLen := len(helpText)
	leftPad := (rightWidth - helpVisibleLen) / 2
	rightPad := rightWidth - helpVisibleLen - leftPad
	paddedHelp := strings.Repeat(" ", leftPad) + styledHelp + strings.Repeat(" ", rightPad)

	rightContent := lipgloss.JoinVertical(lipgloss.Left,
		strings.Repeat(" ", rightWidth),
		styledCmd,
		strings.Repeat(" ", rightWidth),
		paddedHelp,
	)

	// Right box
	rightBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGreen).
		Padding(0, 1)

	rightBox := rightBoxStyle.Render(rightContent)

	// Join with explicit spacer
	spacer := "  "
	installLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, spacer, rightBox)

	// Center the installation layout
	centeredInstallLayout := lipgloss.NewStyle().
		Width(boxWidth - 4). // Account for box padding
		Align(lipgloss.Center).
		Render(installLayout)

	content.WriteString("\n")
	content.WriteString(centeredInstallLayout)
	content.WriteString("\n\n")

	// Help
	help := styles.HelpStyle.Render("j/k or tab: cycle methods • enter/space: copy command • esc/b: back • ?: help • q: quit")
	content.WriteString(help)

	// Wrap in a box with dynamic width
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPinkLight).
		Padding(3, 2).
		Width(boxWidth)

	box := boxStyle.Render(content.String())

	// Center everything
	doc := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(box)

	return doc
}

// formatStringArray converts an array of interfaces to comma-separated string
func formatStringArray(arr []any) string {
	var strs []string
	for _, item := range arr {
		if str, ok := item.(string); ok {
			strs = append(strs, str)
		}
	}
	return strings.Join(strs, ", ")
}

// formatMaintainers formats maintainer information
func formatMaintainers(maintainers []any) string {
	var strs []string
	for _, item := range maintainers {
		switch v := item.(type) {
		case string:
			strs = append(strs, v)
		case map[string]any:
			// Try to extract name, email, or github
			if name, ok := v["name"].(string); ok {
				if email, ok := v["email"].(string); ok {
					strs = append(strs, fmt.Sprintf("%s <%s>", name, email))
				} else if github, ok := v["github"].(string); ok {
					strs = append(strs, fmt.Sprintf("%s (@%s)", name, github))
				} else {
					strs = append(strs, name)
				}
			} else if github, ok := v["github"].(string); ok {
				strs = append(strs, "@"+github)
			}
		}
	}
	return strings.Join(strs, ", ")
}

// wrapText wraps text at specified width
func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		wordLen := len(word)
		if lineLen+wordLen+1 > width && lineLen > 0 {
			result.WriteString("\n")
			lineLen = 0
		}
		if lineLen > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += wordLen

		if i < len(words)-1 && strings.HasSuffix(word, ",") {
			// Natural break point after commas
			if lineLen > width-20 {
				result.WriteString("\n")
				lineLen = 0
			}
		}
	}

	return result.String()
}

// getCurrentPlatform returns the current platform in NixOS format (e.g., "x86_64-darwin")
func getCurrentPlatform() string {
	arch := runtime.GOARCH
	os := runtime.GOOS

	// Map Go arch to Nix arch
	nixArch := arch
	switch arch {
	case "amd64":
		nixArch = "x86_64"
	case "386":
		nixArch = "i686"
	case "arm64":
		nixArch = "aarch64"
	case "arm":
		nixArch = "armv7l"
	}

	// Map Go OS to Nix OS
	nixOS := os
	switch os {
	case "darwin":
		nixOS = "darwin"
	case "linux":
		nixOS = "linux"
	case "windows":
		nixOS = "windows"
	case "freebsd":
		nixOS = "freebsd"
	case "openbsd":
		nixOS = "openbsd"
	case "netbsd":
		nixOS = "netbsd"
	}

	return fmt.Sprintf("%s-%s", nixArch, nixOS)
}

// isPlatformSupported checks if the current platform is in the supported platforms list
func isPlatformSupported(platforms []any, currentPlatform string) bool {
	for _, platform := range platforms {
		if platformStr, ok := platform.(string); ok {
			if platformStr == currentPlatform {
				return true
			}
		}
	}
	return false
}

// renderHelpOverlay renders the help overlay modal
func (m Model) renderHelpOverlay() string {
	helpTitle := lipgloss.NewStyle().
		Foreground(styles.ColorPurple).
		Bold(true).
		Align(lipgloss.Center).
		Render("⌨  KEYBINDINGS REFERENCE")

	// Insert Mode section
	insertModeTitle := lipgloss.NewStyle().
		Foreground(styles.ColorGreen).
		Bold(true).
		Render("INSERT MODE")

	insertModeKeys := []string{
		"Type           → Search for packages",
		"↑ / ↓          → Navigate results",
		"Enter          → Switch to Normal mode",
		"Esc            → Switch to Normal mode",
		"q / Ctrl+C     → Quit application",
	}

	// Normal Mode section
	normalModeTitle := lipgloss.NewStyle().
		Foreground(styles.ColorBlue).
		Bold(true).
		Render("NORMAL MODE")

	normalModeKeys := []string{
		"i              → Switch to Insert mode",
		"j / k          → Move down / up",
		"g / G          → Jump to top / bottom",
		"Enter / Space  → View package details",
		"q / Ctrl+C     → Quit application",
	}

	// Detail View section
	detailModeTitle := lipgloss.NewStyle().
		Foreground(styles.ColorTeal).
		Bold(true).
		Render("DETAIL VIEW")

	detailModeKeys := []string{
		"j / k           → Cycle install methods (down/up)",
		"Tab / Shift+Tab → Cycle install methods (forward/back)",
		"Enter / Space   → Copy selected command",
		"Esc / b         → Back to search",
		"q / Ctrl+C      → Quit application",
	}

	// Global keys
	globalTitle := lipgloss.NewStyle().
		Foreground(styles.ColorYellow).
		Bold(true).
		Render("GLOBAL")

	globalKeys := []string{
		"?              → Toggle this help",
		"q / Ctrl+C     → Quit application",
	}

	// Build content
	var content strings.Builder
	content.WriteString(helpTitle + "\n\n")

	content.WriteString(insertModeTitle + "\n")
	for _, key := range insertModeKeys {
		content.WriteString("  " + lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(key) + "\n")
	}
	content.WriteString("\n")

	content.WriteString(normalModeTitle + "\n")
	for _, key := range normalModeKeys {
		content.WriteString("  " + lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(key) + "\n")
	}
	content.WriteString("\n")

	content.WriteString(detailModeTitle + "\n")
	for _, key := range detailModeKeys {
		content.WriteString("  " + lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(key) + "\n")
	}
	content.WriteString("\n")

	content.WriteString(globalTitle + "\n")
	for _, key := range globalKeys {
		content.WriteString("  " + lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(key) + "\n")
	}

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(styles.ColorGray).
		Italic(true).
		Align(lipgloss.Center).
		Render("\nPress ? or Esc to close")

	content.WriteString("\n" + footer)

	// Create box
	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPurple).
		Padding(1, 2).
		Width(60).
		Render(content.String())

	// Center on screen
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(helpBox)
}

// renderTabMessageOverlay renders the "under development" message for non-Nixpkgs tabs
func (m Model) renderTabMessageOverlay() string {
	tabNames := []string{"Nixpkgs", "Home Manager", "Pacman"}
	selectedTabName := tabNames[m.selectedTab]

	title := lipgloss.NewStyle().
		Foreground(styles.ColorYellow).
		Bold(true).
		Align(lipgloss.Center).
		Render("⚠  UNDER DEVELOPMENT")

	message := lipgloss.NewStyle().
		Foreground(styles.ColorWhite).
		Align(lipgloss.Center).
		Width(50).
		Render(fmt.Sprintf("The %s package source is currently under development and not yet available.", selectedTabName))

	info := lipgloss.NewStyle().
		Foreground(styles.ColorGray).
		Italic(true).
		Align(lipgloss.Center).
		Render("This feature will be implemented in a future release.")

	footer := lipgloss.NewStyle().
		Foreground(styles.ColorGray).
		Italic(true).
		Align(lipgloss.Center).
		Render("\nPress Enter or Esc to continue")

	var content strings.Builder
	content.WriteString(title + "\n\n")
	content.WriteString(message + "\n\n")
	content.WriteString(info)
	content.WriteString(footer)

	// Create box
	messageBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorYellow).
		Padding(2, 4).
		Width(60).
		Render(content.String())

	// Center on screen
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(messageBox)
}
