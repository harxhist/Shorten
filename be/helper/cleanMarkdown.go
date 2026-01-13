package helper

import (
	"regexp"
	"strings"
)

func ExtractTextFromMarkdown(markdownText string) string {
	if markdownText == "" {
		return ""
	}

	text := markdownText

	// Define regex patterns
	patterns := []struct {
		regex   string
		replace string
	}{
		// Remove code blocks (fenced) - with fixed syntax
		// {regex: "(?s)```.*?```", replace: ""},

		// // Remove indented code blocks
		// {regex: `(?m)^(?:\s{4}|\t).*(?:\n|$)`, replace: ""},

		// // Remove inline code
		// {regex: "`[^`]*`", replace: ""},

		// // Remove headers
		// {regex: `(?m)^#{1,6}\s.*$`, replace: ""},

		// // Remove emphasis (bold and italic)
		// {regex: `[*_]{1,2}([^*_]+)[*_]{1,2}`, replace: "$1"},

		// // Remove links but keep link text
		// {regex: `\[([^\]]+)\]\([^\)]+\)`, replace: "$1"},

		// // Remove images
		// {regex: `!\[([^\]]*)\]\([^\)]+\)`, replace: ""},

		// // Remove horizontal rules
		// {regex: `(?m)^[-*_]{3,}\s*$`, replace: ""},

		// // Remove blockquotes
		// {regex: `(?m)^>\s*(.*)`, replace: "$1"},

		// // Remove HTML tags
		// {regex: `<[^>]+>`, replace: ""},

		// // Remove list markers
		// {regex: `(?m)^[\s]*[-+*]\s+`, replace: ""},
		// {regex: `(?m)^[\s]*\d+\.\s+`, replace: ""},

		// Remove all symbols except for essential english language one's
		{regex: `[^a-zA-Z0-9.!:;%$@&"?'\-,\s]`, replace: ""},
	}

	// Apply all regex patterns
	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern.regex)
		text = regex.ReplaceAllString(text, pattern.replace)
	}

	// Clean up extra whitespace and remove '\' and '\n'
	text = strings.ReplaceAll(text, "\n", ". ")
	text = strings.ReplaceAll(text, "\\", "")
	text = strings.TrimSpace(text)

	return text
}