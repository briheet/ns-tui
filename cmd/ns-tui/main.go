package main

import (
	"fmt"
	"os"

	"ns-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	m := ui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
