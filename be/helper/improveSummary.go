package helper

import (
	"strings"
)

// RemoveSummaryLine removes the first line if it starts with "Here is a summary" and contains ":"
func ImproveSummary(inputText string) string {
	if inputText == "" {
		return inputText
	}
	// Split the text into lines
	lines := strings.Split(inputText, "\n")
	if len(lines) == 0 {
		return inputText
	}

	firstLine := lines[0]
	if strings.HasPrefix(firstLine, "Here is a summary") && strings.Contains(firstLine, ":") {
		return strings.TrimLeft(strings.Join(lines[1:], "\n"), "\n")
	}

	return inputText
}
