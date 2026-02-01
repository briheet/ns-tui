package ui

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/briheet/ns-tui/internal/models"
	"github.com/briheet/ns-tui/internal/styles"

	"github.com/charmbracelet/lipgloss"
)

// centerText centers a string within the model's width.
func (m Model) centerText(s string) string {
	return lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(s)
}

// View renders the UI
func (m Model) View() string {
	// Show help overlay
	if m.showHelp {
		return m.renderHelpOverlay()
	}

	// Show HM fetch prompt overlay
	if m.showHMPrompt {
		return m.renderHMPromptOverlay()
	}

	// Show detailed NixOS option view
	if m.mode == models.DetailMode && m.selectedNixOSOption != nil {
		return m.renderNixOSDetailView()
	}

	// Show detailed HM option view
	if m.mode == models.DetailMode && m.selectedHMOption != nil {
		return m.renderHMDetailView()
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

	// HEADER SECTION — use cached banner if width hasn't changed
	var headerContent string
	if m.cache.bannerWidth == m.width && m.cache.banner != "" {
		headerContent = m.cache.banner
	} else {
		headerContent = m.buildBannerHeader()
	}

	// Dynamic parts of the header (depend on state beyond just width)
	var header strings.Builder
	header.WriteString(headerContent)

	// Tabs
	tabNames := []string{"Nixpkgs", "Home Manager", "NixOS Options"}
	var tabParts []string

	for i, name := range tabNames {
		if i > 0 {
			tabParts = append(tabParts, styles.TabSeparatorStyle.Render("│"))
		}
		if i == m.selectedTab {
			tabParts = append(tabParts, styles.ActiveTabStyle.Render(name))
		} else {
			tabParts = append(tabParts, styles.InactiveTabStyle.Render(name))
		}
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabParts...)
	header.WriteString(m.centerText(tabs))
	header.WriteString("\n\n")

	// Search box with mode indicator
	var modeIndicator strings.Builder
	if m.mode == models.InsertMode {
		modeIndicator.WriteString(styles.InsertModeStyle.Render("-- INSERT --"))
	} else {
		modeIndicator.WriteString(styles.NormalModeStyle.Render("-- NORMAL --"))
	}

	if m.selectedTab == 0 && len(m.packages) > 0 {
		modeIndicator.WriteString(styles.PositionStyle.Render(fmt.Sprintf(" [%d/%d]", m.cursor+1, len(m.packages))))
	} else if m.selectedTab == 1 && len(m.hmSearchResults) > 0 {
		modeIndicator.WriteString(styles.PositionStyle.Render(fmt.Sprintf(" [%d/%d]", m.hmCursor+1, len(m.hmSearchResults))))
	} else if m.selectedTab == 2 && len(m.nixosSearchResults) > 0 {
		modeIndicator.WriteString(styles.PositionStyle.Render(fmt.Sprintf(" [%d/%d]", m.nixosCursor+1, len(m.nixosSearchResults))))
	}

	searchBox := styles.SearchBoxStyle.Render(m.textInput.View())
	header.WriteString(m.centerText(searchBox))
	header.WriteString("\n")
	header.WriteString(m.centerText(modeIndicator.String()))
	header.WriteString("\n")

	separator := styles.SeparatorStyle.Render(strings.Repeat("─", min(m.width, 80)))
	header.WriteString(m.centerText(separator))
	header.WriteString("\n\n")

	headerStr := header.String()
	headerLines := strings.Count(headerStr, "\n")

	// RESULTS SECTION — pass remaining height so results don't overflow
	remainingLines := m.height - headerLines - footerHeight - 1
	var resultsContent string
	switch m.selectedTab {
	case 1:
		resultsContent = m.renderHMResults(remainingLines)
	case 2:
		resultsContent = m.renderNixOSResults(remainingLines)
	default:
		resultsContent = m.renderResults(remainingLines)
	}

	// FOOTER SECTION (fixed at bottom)
	var footer strings.Builder

	// Toast (if visible)
	if m.toastVisible {
		toast := styles.ToastStyle.Render(m.toastMessage)
		footer.WriteString(m.centerText(toast))
		footer.WriteString("\n")
	}

	// Help bar — use cached if width matches
	var helpRendered string
	if m.cache.helpWidth == m.width {
		if m.mode == models.InsertMode {
			helpRendered = m.cache.helpInsert
		} else {
			helpRendered = m.cache.helpNormal
		}
	}
	if helpRendered == "" {
		var helpText string
		if m.mode == models.InsertMode {
			helpText = "esc: normal mode • ↑/↓: navigate • tab: switch source • ?: help • q: quit"
		} else {
			helpText = "i: insert mode • j/k: navigate • enter/space: details • g/G: top/bottom • tab: switch source • ?: help • q: quit"
		}
		helpRendered = m.centerText(styles.HelpStyle.Render(helpText))
	}
	footer.WriteString(helpRendered)

	footerContent := footer.String()

	// COMBINE ALL SECTIONS
	var result strings.Builder
	result.WriteString(headerStr)
	result.WriteString(resultsContent)

	// Fill remaining space to push footer to bottom
	mainLines := strings.Count(result.String(), "\n")
	availableHeight := m.height - footerHeight - 1
	if mainLines < availableHeight {
		result.WriteString(strings.Repeat("\n", availableHeight-mainLines))
	}

	result.WriteString("\n")
	result.WriteString(footerContent)
	return result.String()
}

// buildBannerHeader builds the static banner portion of the header.
func (m Model) buildBannerHeader() string {
	var header strings.Builder
	header.WriteString("\n\n\n")

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
		styledLine := styles.BannerStyle.Render(line)
		header.WriteString(m.centerText(styledLine))
		header.WriteString("\n")
	}

	dotsLine := styles.BannerGrayStyle.Render("..::::..:::......:::::::::::::::..::::::.......:::.....::")
	header.WriteString(m.centerText(dotsLine))
	header.WriteString("\n\n")

	var subtitleText string
	switch m.selectedTab {
	case 0:
		subtitleText = "Real-time package discovery with fuzzy search"
	case 1:
		subtitleText = "Search Home Manager configuration options"
	case 2:
		subtitleText = "Search NixOS configuration options"
	default:
		subtitleText = "Real-time package discovery with fuzzy search"
	}
	subtitle := styles.SubtitleStyle.Render(subtitleText)
	header.WriteString(m.centerText(subtitle))
	header.WriteString("\n\n")

	return header.String()
}

// renderResults renders the search results within the given line budget
func (m Model) renderResults(availHeight int) string {
	var content strings.Builder

	// Loading indicator with spinner
	if m.loading {
		loading := fmt.Sprintf("%s Searching...", m.spinner.View())
		loadingStyled := styles.LoadingStyle.Render(loading)
		content.WriteString(m.centerText(loadingStyled))
		content.WriteString("\n")
		return content.String()
	}

	// Error message
	if m.err != nil {
		errorMsg := styles.ErrorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err))
		content.WriteString(m.centerText(errorMsg))
		content.WriteString("\n")
		return content.String()
	}

	// Results
	if len(m.packages) > 0 {
		// Reserve ~4 lines for count header and scroll indicators, each item ~4 lines
		maxVisible := (availHeight - 4) / 4
		if maxVisible < 3 {
			maxVisible = 3
		}

		visibleCount := min(maxVisible, len(m.packages))
		count := styles.CountStyle.Render(fmt.Sprintf("📦 %d packages (showing %d)", len(m.packages), visibleCount))
		content.WriteString(m.centerText(count))
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
			scrollUp := styles.ScrollIndicatorStyle.Render("⬆ More above")
			content.WriteString(m.centerText(scrollUp))
			content.WriteString("\n")
		}
		if end < len(m.packages) {
			scrollDown := styles.ScrollIndicatorStyle.Render(fmt.Sprintf("⬇ %d more below", len(m.packages)-end))
			content.WriteString(m.centerText(scrollDown))
			content.WriteString("\n")
		}
	} else if !m.loading && m.lastQuery != "" {
		noResults := styles.NoResultsStyle.Render("No packages found. Try a different search term.")
		content.WriteString(m.centerText(noResults))
		content.WriteString("\n")
	} else if m.lastQuery == "" && !m.loading {
		hint := styles.HintStyle.Render("Type to search for NixOS packages...")
		content.WriteString(m.centerText(hint))
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
	version := styles.VersionStyle.Render("v" + pkg.Version)
	desc := strings.Join(strings.Fields(pkg.Description), " ")
	maxDescLen := 70
	if m.width > 100 {
		maxDescLen = 100
	}
	if len(desc) > maxDescLen {
		desc = desc[:maxDescLen-3] + "..."
	}
	desc = styles.DescriptionStyle.Render(desc)

	// Build line with strings.Builder
	var b strings.Builder
	b.WriteString(cursor)
	b.WriteString(name)
	b.WriteByte(' ')
	b.WriteString(version)
	b.WriteString("\n     ")
	b.WriteString(desc)

	if len(pkg.Programs) > 0 {
		var programsList []string
		maxPrograms := 5
		for i, prog := range pkg.Programs {
			if i >= maxPrograms {
				programsList = append(programsList, fmt.Sprintf("+%d more", len(pkg.Programs)-maxPrograms))
				break
			}
			programsList = append(programsList, prog)
		}
		programsText := styles.ProgramsStyle.Render("📦 " + strings.Join(programsList, ", "))
		b.WriteString("\n     ")
		b.WriteString(programsText)
	}

	line := b.String()

	var renderedLine string
	if m.cursor == index {
		renderedLine = styles.SelectedItemStyle.Render(line)
	} else {
		renderedLine = styles.ResultItemStyle.Render(line)
	}

	return m.centerText(renderedLine) + "\n"
}

// renderHMResults renders the Home Manager search results within the given line budget
func (m Model) renderHMResults(availHeight int) string {
	var content strings.Builder

	// Loading indicator with spinner (during fetch)
	if m.hmLoading {
		loading := fmt.Sprintf("%s Fetching Home Manager options...", m.spinner.View())
		loadingStyled := styles.LoadingStyle.Render(loading)
		content.WriteString(m.centerText(loadingStyled))
		content.WriteString("\n")
		return content.String()
	}

	// Loading indicator for search
	if m.loading {
		loading := fmt.Sprintf("%s Searching...", m.spinner.View())
		loadingStyled := styles.LoadingStyle.Render(loading)
		content.WriteString(m.centerText(loadingStyled))
		content.WriteString("\n")
		return content.String()
	}

	// Error message
	if m.hmErr != nil {
		errorMsg := styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.hmErr))
		content.WriteString(m.centerText(errorMsg))
		content.WriteString("\n")
		return content.String()
	}

	// Not loaded yet (waiting for user action)
	if !m.hmLoaded {
		hint := styles.HintGrayStyle.Render("Home Manager options not loaded yet.")
		content.WriteString(m.centerText(hint))
		content.WriteString("\n")
		return content.String()
	}

	// Results
	if len(m.hmSearchResults) > 0 {
		// Reserve ~4 lines for count header and scroll indicators, each HM item ~3 lines
		maxVisible := (availHeight - 4) / 3
		if maxVisible < 3 {
			maxVisible = 3
		}

		visibleCount := min(maxVisible, len(m.hmSearchResults))
		count := styles.CountStyle.Render(fmt.Sprintf("Found %d options (showing %d)", len(m.hmSearchResults), visibleCount))
		content.WriteString(m.centerText(count))
		content.WriteString("\n\n")

		// Ensure cursor is in view
		if m.hmCursor < m.hmScrollOffset {
			m.hmScrollOffset = m.hmCursor
		}
		if m.hmCursor >= m.hmScrollOffset+maxVisible {
			m.hmScrollOffset = m.hmCursor - maxVisible + 1
		}

		start := m.hmScrollOffset
		end := min(m.hmScrollOffset+maxVisible, len(m.hmSearchResults))

		for i := start; i < end; i++ {
			content.WriteString(m.renderHMOptionItem(i))
		}

		// Scroll indicators
		if m.hmScrollOffset > 0 {
			scrollUp := styles.ScrollIndicatorStyle.Render("More above")
			content.WriteString(m.centerText(scrollUp))
			content.WriteString("\n")
		}
		if end < len(m.hmSearchResults) {
			scrollDown := styles.ScrollIndicatorStyle.Render(fmt.Sprintf("%d more below", len(m.hmSearchResults)-end))
			content.WriteString(m.centerText(scrollDown))
			content.WriteString("\n")
		}
	} else if m.hmLastQuery != "" {
		noResults := styles.NoResultsStyle.Render("No options found. Try a different search term.")
		content.WriteString(m.centerText(noResults))
		content.WriteString("\n")
	} else {
		hint := styles.HintStyle.Render("Type to search Home Manager options...")
		content.WriteString(m.centerText(hint))
		content.WriteString("\n")
	}

	return content.String()
}

// renderHMOptionItem renders a single Home Manager option item
func (m Model) renderHMOptionItem(index int) string {
	opt := m.hmSearchResults[index]

	cursor := "  "
	if m.hmCursor == index {
		cursor = "▶ "
	}

	name := styles.PackageNameStyle.Render(opt.Name)
	typeStr := styles.VersionStyle.Render(opt.Type)

	// Collapse newlines and trim to single line
	desc := strings.Join(strings.Fields(opt.Description), " ")
	maxDescLen := 70
	if m.width > 100 {
		maxDescLen = 100
	}
	if len(desc) > maxDescLen {
		desc = desc[:maxDescLen-3] + "..."
	}
	desc = styles.DescriptionStyle.Render(desc)

	var b strings.Builder
	b.WriteString(cursor)
	b.WriteString(name)
	b.WriteString("  ")
	b.WriteString(typeStr)
	b.WriteString("\n     ")
	b.WriteString(desc)
	line := b.String()

	var renderedLine string
	if m.hmCursor == index {
		renderedLine = styles.SelectedItemStyle.Render(line)
	} else {
		renderedLine = styles.ResultItemStyle.Render(line)
	}

	return m.centerText(renderedLine) + "\n"
}

// renderHMPromptOverlay renders the Home Manager fetch prompt modal
func (m Model) renderHMPromptOverlay() string {
	title := styles.OverlayTitleStyle.Render("HOME MANAGER OPTIONS")

	message := styles.OverlayMessageStyle.Width(50).
		Render("Home Manager options have not been fetched yet. This will run `nix build` to download and cache the options JSON locally.")

	info := styles.OverlayInfoStyle.Render("This may take a minute on first run.")

	// Yes / No buttons
	yesStyle := lipgloss.NewStyle().Padding(0, 3)
	noStyle := lipgloss.NewStyle().Padding(0, 3)

	if m.hmPromptSelection == 0 {
		yesStyle = yesStyle.Background(styles.ColorGreen).Foreground(styles.ColorBg).Bold(true)
		noStyle = noStyle.Foreground(styles.ColorGray)
	} else {
		yesStyle = yesStyle.Foreground(styles.ColorGray)
		noStyle = noStyle.Background(styles.ColorRed).Foreground(styles.ColorBg).Bold(true)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Top,
		yesStyle.Render("Yes, fetch"),
		"   ",
		noStyle.Render("No, go back"),
	)

	footer := styles.OverlayInfoStyle.Render("\nj/k to toggle • Enter to confirm • Esc to cancel")

	var content strings.Builder
	content.WriteString(title + "\n\n")
	content.WriteString(message + "\n\n")
	content.WriteString(info + "\n\n")
	centeredButtons := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(buttons)
	content.WriteString(centeredButtons)
	content.WriteString("\n" + footer)

	// Create box
	promptBox := styles.OverlayBoxStyle.Render(content.String())

	// Center on screen
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(promptBox)
}

// renderHMDetailView renders the detailed Home Manager option view
func (m Model) renderHMDetailView() string {
	var content strings.Builder
	opt := m.selectedHMOption

	// Calculate responsive box width (same as Nixpkgs detail)
	boxWidth := int(float64(m.width) * 0.9)
	if boxWidth > 160 {
		boxWidth = 160
	}
	centerStyle := lipgloss.NewStyle().Width(boxWidth - 4).Align(lipgloss.Center)

	// 1. Breadcrumb from loc
	breadcrumb := buildBreadcrumb(opt.Loc)
	content.WriteString(centerStyle.Render(breadcrumb))
	content.WriteString("\n\n")

	// 2. Type
	typeLine := styles.DetailLabelStyle.Render("Type: ") + styles.DetailValueStyle.Render(opt.Type)
	content.WriteString(typeLine)
	content.WriteString("\n")

	// 3. Read-only (only show if true)
	if opt.ReadOnly {
		roLine := styles.DetailLabelStyle.Render("Read Only: ") +
			styles.ReadOnlyStyle.Render("Yes")
		content.WriteString(roLine)
		content.WriteString("\n")
	}
	content.WriteString("\n")

	// 4. Description
	if opt.Description != "" {
		content.WriteString(styles.DetailLabelStyle.Render("Description:"))
		content.WriteString("\n")
		wrapped := wrapText(strings.TrimSpace(opt.Description), boxWidth-8)
		content.WriteString(styles.DetailValueStyle.Render(wrapped))
		content.WriteString("\n\n")
	}

	// 5. Default value
	if opt.Default != nil {
		defLine := styles.DetailLabelStyle.Render("Default: ") + styles.DetailValueStyle.Render(*opt.Default)
		content.WriteString(defLine)
		content.WriteString("\n\n")
	}

	// 6. Example value
	if opt.Example != nil {
		exLine := styles.DetailLabelStyle.Render("Example: ") + styles.DetailValueStyle.Render(*opt.Example)
		content.WriteString(exLine)
		content.WriteString("\n\n")
	}

	// 7. Declarations / source links
	if len(opt.Declarations) > 0 {
		content.WriteString(styles.DetailLabelStyle.Render("Declared in:"))
		content.WriteString("\n")
		for _, decl := range opt.Declarations {
			content.WriteString("  " + styles.DetailValueStyle.Render(decl.Name))
			if decl.URL != "" {
				content.WriteString("\n  " + styles.URLStyle.Render(decl.URL))
			}
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	// 8. Related options section
	if len(m.hmRelatedOptions) > 0 {
		parentPath := strings.Join(opt.Loc[:len(opt.Loc)-1], ".")
		separatorLabel := fmt.Sprintf(" Related Options (%s.*) ", parentPath)
		lineLen := boxWidth - 8 - len(separatorLabel)
		leftLine := strings.Repeat("─", 2)
		rightLine := ""
		if lineLen > 2 {
			rightLine = strings.Repeat("─", lineLen-2)
		}
		sep := styles.SeparatorStyle.Render(leftLine + separatorLabel + rightLine)
		content.WriteString(sep)
		content.WriteString("\n\n")

		maxVisible := 8
		if m.height > 30 {
			maxVisible = 12
		}

		// Ensure cursor is in view
		hmRelScroll := m.hmRelatedScrollOffset
		if m.hmRelatedCursor < hmRelScroll {
			hmRelScroll = m.hmRelatedCursor
		}
		if m.hmRelatedCursor >= hmRelScroll+maxVisible {
			hmRelScroll = m.hmRelatedCursor - maxVisible + 1
		}

		start := hmRelScroll
		end := min(start+maxVisible, len(m.hmRelatedOptions))

		relCenterStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(boxWidth - 4)

		for i := start; i < end; i++ {
			rel := m.hmRelatedOptions[i]
			cur := "  "
			if m.hmRelatedCursor == i {
				cur = "▶ "
			}

			nameStyled := styles.PackageNameStyle.Render(rel.Name)
			typeStyled := styles.VersionStyle.Render(rel.Type)

			// Truncate description for the related item
			desc := rel.Description
			maxDescLen := 60
			if m.width > 100 {
				maxDescLen = 90
			}
			if len(desc) > maxDescLen {
				desc = desc[:maxDescLen-3] + "..."
			}
			descStyled := styles.DescriptionStyle.Render(desc)

			var b strings.Builder
			b.WriteString(cur)
			b.WriteString(nameStyled)
			b.WriteString("  ")
			b.WriteString(typeStyled)
			b.WriteString("\n     ")
			b.WriteString(descStyled)
			line := b.String()

			if m.hmRelatedCursor == i {
				line = styles.SelectedItemStyle.Render(line)
			} else {
				line = styles.ResultItemStyle.Render(line)
			}
			content.WriteString(relCenterStyle.Render(line))
			content.WriteString("\n\n")
		}

		// Scroll indicators
		if start > 0 {
			content.WriteString(styles.ScrollIndicatorStyle.Render("  More above"))
			content.WriteString("\n")
		}
		if end < len(m.hmRelatedOptions) {
			content.WriteString(styles.ScrollIndicatorStyle.Render(fmt.Sprintf("  %d more below", len(m.hmRelatedOptions)-end)))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	} else {
		content.WriteString(styles.HintGrayStyle.Render("No related options found."))
		content.WriteString("\n\n")
	}

	// 9. Help bar
	helpText := "j/k: navigate related • enter/space: view option • esc/b: back • ?: help • q: quit"
	help := styles.HelpStyle.Render(helpText)
	content.WriteString(help)

	// 10. Wrap in box — use most of the screen height
	boxHeight := m.height - 12
	if boxHeight < 12 {
		boxHeight = 12
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPinkLight).
		Padding(1, 2).
		Width(boxWidth).
		Height(boxHeight)

	box := boxStyle.Render(content.String())

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(box)
}

// buildBreadcrumb renders a styled breadcrumb from loc segments
func buildBreadcrumb(loc []string) string {
	if len(loc) == 0 {
		return ""
	}

	var parts []string
	for i, seg := range loc {
		parts = append(parts, styles.BreadcrumbSegmentStyle.Render(seg))
		if i < len(loc)-1 {
			parts = append(parts, styles.BreadcrumbSepStyle.Render(" > "))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
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
	if len(pkg.HomepageLinks) > 0 {
		content.WriteString(styles.DetailLabelStyle.Render("Homepage: "))
		homeLinkList := strings.Join(pkg.HomepageLinks, ", ")
		wrappedHomeLinks := wrapText(homeLinkList, boxWidth-8)
		content.WriteString(styles.DetailValueStyle.Render(wrappedHomeLinks))
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
		programsList := strings.Join(pkg.Programs, ", ")
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
			content.WriteString(styles.PlatformSupportedStyle.Render("✓ Supported on your platform"))
			content.WriteString(styles.DetailValueStyle.Render(fmt.Sprintf(" (%s)", currentPlatform)))
		} else {
			content.WriteString(styles.PlatformUnsupportedStyle.Render("✗ Not supported on your platform"))
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

	var menuItems []string
	for i, name := range methodNames {
		var line string
		if i == m.selectedInstallMethod {
			indicator := styles.IndicatorStyle.Render("→ ")
			methodName := styles.SelectedMethodStyle.Render(name)
			line = indicator + methodName
			visibleLen := 2 + len(name)
			if visibleLen < leftWidth {
				line = line + strings.Repeat(" ", leftWidth-visibleLen)
			}
		} else {
			methodName := styles.NormalMethodStyle.Render(name)
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
	leftBox := styles.LeftBoxStyle.Render(leftContent)

	// Right box content - styled command
	selectedCmd := commands[m.selectedInstallMethod]

	// Style the command with green color and center it
	styledCmd := styles.CmdStyle.Render(selectedCmd)

	// Center the command
	visibleCmdLen := len(selectedCmd)
	if visibleCmdLen < rightWidth {
		cmdLeftPad := (rightWidth - visibleCmdLen) / 2
		cmdRightPad := rightWidth - visibleCmdLen - cmdLeftPad
		styledCmd = strings.Repeat(" ", cmdLeftPad) + styledCmd + strings.Repeat(" ", cmdRightPad)
	} else if visibleCmdLen > rightWidth {
		selectedCmd = selectedCmd[:rightWidth]
		styledCmd = styles.CmdStyle.Render(selectedCmd)
	}

	// Style help text - bold and visible
	helpStyle := styles.CopyHelpStyle
	helpText := "Press Enter or Spacebar to copy"

	// Change text when copied successfully
	if m.toastVisible {
		helpStyle = styles.CopiedHelpStyle
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
	rightBox := styles.RightBoxStyle.Render(rightContent)

	// Join with explicit spacer
	spacer := "  "
	installLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, spacer, rightBox)

	// Center the installation layout
	centeredInstallLayout := lipgloss.NewStyle().
		Width(boxWidth - 4).
		Align(lipgloss.Center).
		Render(installLayout)

	content.WriteString("\n")
	content.WriteString(centeredInstallLayout)
	content.WriteString("\n\n")

	// Help
	help := styles.HelpStyle.Render("j/k or tab: cycle methods • enter/space: copy command • esc/b: back • ?: help • q: quit")
	content.WriteString(help)

	// Wrap in a box with dynamic width
	nixBoxHeight := m.height - 12
	if nixBoxHeight < 12 {
		nixBoxHeight = 12
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPinkLight).
		Padding(1, 2).
		Width(boxWidth).
		Height(nixBoxHeight)

	box := boxStyle.Render(content.String())

	// Center everything
	doc := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(box)

	return doc
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

	return nixArch + "-" + nixOS
}

// isPlatformSupported checks if the current platform is in the supported platforms list
func isPlatformSupported(platforms []string, currentPlatform string) bool {
	for _, platform := range platforms {
		if platform == currentPlatform {
			return true
		}
	}
	return false
}

// renderHelpOverlay renders the help overlay modal
func (m Model) renderHelpOverlay() string {
	helpTitle := styles.OverlayTitleStyle.Render("⌨  KEYBINDINGS REFERENCE")

	// Insert Mode section
	insertModeTitle := styles.PlatformSupportedStyle.Render("INSERT MODE")

	insertModeKeys := []string{
		"Type           → Search for packages",
		"↑ / ↓          → Navigate results",
		"Enter          → Switch to Normal mode",
		"Esc            → Switch to Normal mode",
		"q / Ctrl+C     → Quit application",
	}

	// Normal Mode section
	normalModeTitle := styles.URLStyle.Render("NORMAL MODE")

	normalModeKeys := []string{
		"i              → Switch to Insert mode",
		"j / k          → Move down / up",
		"g / G          → Jump to top / bottom",
		"Enter / Space  → View package details",
		"q / Ctrl+C     → Quit application",
	}

	// Detail View section
	detailModeTitle := styles.CopyHelpStyle.Render("DETAIL VIEW")

	detailModeKeys := []string{
		"j / k           → Cycle install methods (down/up)",
		"Tab / Shift+Tab → Cycle install methods (forward/back)",
		"Enter / Space   → Copy selected command",
		"Esc / b         → Back to search",
		"q / Ctrl+C      → Quit application",
	}

	// Global keys
	globalTitle := styles.ReadOnlyStyle.Render("GLOBAL")

	globalKeys := []string{
		"?              → Toggle this help",
		"q / Ctrl+C     → Quit application",
	}

	// Build content
	var content strings.Builder
	content.WriteString(helpTitle + "\n\n")

	content.WriteString(insertModeTitle + "\n")
	for _, key := range insertModeKeys {
		content.WriteString("  " + styles.WhiteTextStyle.Render(key) + "\n")
	}
	content.WriteString("\n")

	content.WriteString(normalModeTitle + "\n")
	for _, key := range normalModeKeys {
		content.WriteString("  " + styles.WhiteTextStyle.Render(key) + "\n")
	}
	content.WriteString("\n")

	content.WriteString(detailModeTitle + "\n")
	for _, key := range detailModeKeys {
		content.WriteString("  " + styles.WhiteTextStyle.Render(key) + "\n")
	}
	content.WriteString("\n")

	content.WriteString(globalTitle + "\n")
	for _, key := range globalKeys {
		content.WriteString("  " + styles.WhiteTextStyle.Render(key) + "\n")
	}

	// Footer
	footer := styles.OverlayInfoStyle.Render("\nPress ? or Esc to close")

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

// renderNixOSResults renders the NixOS options search results within the given line budget
func (m Model) renderNixOSResults(availHeight int) string {
	var content strings.Builder

	// Loading indicator with spinner
	if m.loading {
		loading := fmt.Sprintf("%s Searching...", m.spinner.View())
		loadingStyled := styles.LoadingStyle.Render(loading)
		content.WriteString(m.centerText(loadingStyled))
		content.WriteString("\n")
		return content.String()
	}

	// Error message
	if m.nixosErr != nil {
		errorMsg := styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.nixosErr))
		content.WriteString(m.centerText(errorMsg))
		content.WriteString("\n")
		return content.String()
	}

	// Results
	if len(m.nixosSearchResults) > 0 {
		// Reserve ~4 lines for count header and scroll indicators, each item ~3 lines
		maxVisible := (availHeight - 4) / 3
		if maxVisible < 3 {
			maxVisible = 3
		}

		visibleCount := min(maxVisible, len(m.nixosSearchResults))
		count := styles.CountStyle.Render(fmt.Sprintf("Found %d options (showing %d)", len(m.nixosSearchResults), visibleCount))
		content.WriteString(m.centerText(count))
		content.WriteString("\n\n")

		// Ensure cursor is in view
		if m.nixosCursor < m.nixosScrollOffset {
			m.nixosScrollOffset = m.nixosCursor
		}
		if m.nixosCursor >= m.nixosScrollOffset+maxVisible {
			m.nixosScrollOffset = m.nixosCursor - maxVisible + 1
		}

		start := m.nixosScrollOffset
		end := min(m.nixosScrollOffset+maxVisible, len(m.nixosSearchResults))

		for i := start; i < end; i++ {
			content.WriteString(m.renderNixOSOptionItem(i))
		}

		// Scroll indicators
		if m.nixosScrollOffset > 0 {
			scrollUp := styles.ScrollIndicatorStyle.Render("More above")
			content.WriteString(m.centerText(scrollUp))
			content.WriteString("\n")
		}
		if end < len(m.nixosSearchResults) {
			scrollDown := styles.ScrollIndicatorStyle.Render(fmt.Sprintf("%d more below", len(m.nixosSearchResults)-end))
			content.WriteString(m.centerText(scrollDown))
			content.WriteString("\n")
		}
	} else if m.nixosLastQuery != "" {
		noResults := styles.NoResultsStyle.Render("No options found. Try a different search term.")
		content.WriteString(m.centerText(noResults))
		content.WriteString("\n")
	} else {
		hint := styles.HintStyle.Render("Type to search NixOS configuration options...")
		content.WriteString(m.centerText(hint))
		content.WriteString("\n")
	}

	return content.String()
}

// renderNixOSOptionItem renders a single NixOS option item
func (m Model) renderNixOSOptionItem(index int) string {
	opt := m.nixosSearchResults[index]

	cursor := "  "
	if m.nixosCursor == index {
		cursor = "▶ "
	}

	name := styles.PackageNameStyle.Render(opt.Name)
	typeStr := styles.VersionStyle.Render(opt.Type)

	// Collapse newlines, strip HTML, and trim to single line
	desc := stripHTML(opt.Description)
	desc = strings.Join(strings.Fields(desc), " ")
	maxDescLen := 70
	if m.width > 100 {
		maxDescLen = 100
	}
	if len(desc) > maxDescLen {
		desc = desc[:maxDescLen-3] + "..."
	}
	desc = styles.DescriptionStyle.Render(desc)

	var b strings.Builder
	b.WriteString(cursor)
	b.WriteString(name)
	b.WriteString("  ")
	b.WriteString(typeStr)
	b.WriteString("\n     ")
	b.WriteString(desc)
	line := b.String()

	var renderedLine string
	if m.nixosCursor == index {
		renderedLine = styles.SelectedItemStyle.Render(line)
	} else {
		renderedLine = styles.ResultItemStyle.Render(line)
	}

	return m.centerText(renderedLine) + "\n"
}

// renderNixOSDetailView renders the detailed NixOS option view
func (m Model) renderNixOSDetailView() string {
	var content strings.Builder
	opt := m.selectedNixOSOption

	// Calculate responsive box width (same as other detail views)
	boxWidth := int(float64(m.width) * 0.9)
	if boxWidth > 160 {
		boxWidth = 160
	}
	centerStyle := lipgloss.NewStyle().Width(boxWidth - 4).Align(lipgloss.Center)

	// 1. Breadcrumb from loc
	breadcrumb := buildBreadcrumb(opt.Loc)
	content.WriteString(centerStyle.Render(breadcrumb))
	content.WriteString("\n\n")

	// 2. Type
	typeLine := styles.DetailLabelStyle.Render("Type: ") + styles.DetailValueStyle.Render(opt.Type)
	content.WriteString(typeLine)
	content.WriteString("\n\n")

	// 3. Description
	if opt.Description != "" {
		content.WriteString(styles.DetailLabelStyle.Render("Description:"))
		content.WriteString("\n")
		desc := stripHTML(opt.Description)
		wrapped := wrapText(strings.TrimSpace(desc), boxWidth-8)
		content.WriteString(styles.DetailValueStyle.Render(wrapped))
		content.WriteString("\n\n")
	}

	// 4. Default value
	if opt.Default != nil {
		defLine := styles.DetailLabelStyle.Render("Default: ") + styles.DetailValueStyle.Render(*opt.Default)
		content.WriteString(defLine)
		content.WriteString("\n\n")
	}

	// 5. Example value
	if opt.Example != nil {
		exLine := styles.DetailLabelStyle.Render("Example: ") + styles.DetailValueStyle.Render(*opt.Example)
		content.WriteString(exLine)
		content.WriteString("\n\n")
	}

	// 6. Source file + GitHub link
	if opt.Source != "" {
		content.WriteString(styles.DetailLabelStyle.Render("Source:"))
		content.WriteString("\n")
		content.WriteString("  " + styles.DetailValueStyle.Render(opt.Source))
		content.WriteString("\n")
		ghURL := "https://github.com/NixOS/nixpkgs/blob/nixos-unstable/" + opt.Source
		content.WriteString("  " + styles.URLStyle.Render(ghURL))
		content.WriteString("\n\n")
	}

	// 7. Related options section
	if m.nixosRelatedLoading {
		loading := fmt.Sprintf("%s Loading related options...", m.spinner.View())
		content.WriteString(styles.LoadingStyle.Render(loading))
		content.WriteString("\n\n")
	} else if len(m.nixosRelatedOptions) > 0 {
		parentPath := strings.Join(opt.Loc[:len(opt.Loc)-1], ".")
		separatorLabel := fmt.Sprintf(" Related Options (%s.*) ", parentPath)
		lineLen := boxWidth - 8 - len(separatorLabel)
		leftLine := strings.Repeat("─", 2)
		rightLine := ""
		if lineLen > 2 {
			rightLine = strings.Repeat("─", lineLen-2)
		}
		sep := styles.SeparatorStyle.Render(leftLine + separatorLabel + rightLine)
		content.WriteString(sep)
		content.WriteString("\n\n")

		maxVisible := 8
		if m.height > 30 {
			maxVisible = 12
		}

		// Ensure cursor is in view
		nixRelScroll := m.nixosRelatedScrollOffset
		if m.nixosRelatedCursor < nixRelScroll {
			nixRelScroll = m.nixosRelatedCursor
		}
		if m.nixosRelatedCursor >= nixRelScroll+maxVisible {
			nixRelScroll = m.nixosRelatedCursor - maxVisible + 1
		}

		start := nixRelScroll
		end := min(start+maxVisible, len(m.nixosRelatedOptions))

		relCenterStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(boxWidth - 4)

		for i := start; i < end; i++ {
			rel := m.nixosRelatedOptions[i]
			cur := "  "
			if m.nixosRelatedCursor == i {
				cur = "▶ "
			}

			nameStyled := styles.PackageNameStyle.Render(rel.Name)
			typeStyled := styles.VersionStyle.Render(rel.Type)

			// Truncate description for the related item
			desc := stripHTML(rel.Description)
			maxDescLen := 60
			if m.width > 100 {
				maxDescLen = 90
			}
			if len(desc) > maxDescLen {
				desc = desc[:maxDescLen-3] + "..."
			}
			descStyled := styles.DescriptionStyle.Render(desc)

			var b strings.Builder
			b.WriteString(cur)
			b.WriteString(nameStyled)
			b.WriteString("  ")
			b.WriteString(typeStyled)
			b.WriteString("\n     ")
			b.WriteString(descStyled)
			line := b.String()

			if m.nixosRelatedCursor == i {
				line = styles.SelectedItemStyle.Render(line)
			} else {
				line = styles.ResultItemStyle.Render(line)
			}
			content.WriteString(relCenterStyle.Render(line))
			content.WriteString("\n\n")
		}

		// Scroll indicators
		if start > 0 {
			content.WriteString(styles.ScrollIndicatorStyle.Render("  More above"))
			content.WriteString("\n")
		}
		if end < len(m.nixosRelatedOptions) {
			content.WriteString(styles.ScrollIndicatorStyle.Render(fmt.Sprintf("  %d more below", len(m.nixosRelatedOptions)-end)))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	} else {
		content.WriteString(styles.HintGrayStyle.Render("No related options found."))
		content.WriteString("\n\n")
	}

	// 8. Help bar
	helpText := "j/k: navigate related • enter/space: view option • esc/b: back • ?: help • q: quit"
	help := styles.HelpStyle.Render(helpText)
	content.WriteString(help)

	// 9. Wrap in box
	boxHeight := m.height - 12
	if boxHeight < 12 {
		boxHeight = 12
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPinkLight).
		Padding(1, 2).
		Width(boxWidth).
		Height(boxHeight)

	box := boxStyle.Render(content.String())

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(box)
}
