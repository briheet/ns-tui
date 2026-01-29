package ui

import "github.com/briheet/ns-tui/internal/models"

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

// hmCacheCheckMsg is sent after checking whether the HM cache exists
type hmCacheCheckMsg struct {
	exists  bool
	options []models.HMOption
	err     error
}

// hmFetchResultMsg is sent when the HM options fetch completes
type hmFetchResultMsg struct {
	options []models.HMOption
	err     error
}

// hmSearchResultMsg is sent when an in-memory HM search completes
type hmSearchResultMsg struct {
	results []models.HMOption
}
