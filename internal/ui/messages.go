package ui

import "ns-tui/internal/models"

// searchResultMsg is sent when search results are received
type searchResultMsg struct {
	packages []models.Package
	err      error
}

// clipboardMsg is sent when clipboard operation completes
type clipboardMsg struct {
	success bool
	command string
	err     error
}

// hideToastMsg is sent when the toast should be hidden
type hideToastMsg struct{}
