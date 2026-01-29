package ui

import (
	"regexp"
	"strings"
)

var (
	// HTML tag regex
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
	// HTML entity regex
	htmlEntityRegex = regexp.MustCompile(`&[a-zA-Z]+;|&#[0-9]+;`)
)

// stripHTML removes HTML tags and entities from a string
func stripHTML(html string) string {
	// Remove HTML tags
	text := htmlTagRegex.ReplaceAllString(html, "")

	// Replace common HTML entities
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&nbsp;", " ")

	// Replace any remaining entities with space
	text = htmlEntityRegex.ReplaceAllString(text, " ")

	// Clean up multiple spaces and newlines
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	return text
}
