package helper

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
	"unicode"
)

// ExtractText extracts text from specific HTML tags, ensuring that the final string has no more than 28000 words.
func ExtractText(htmlData string) (string, error) {
	// Parse the HTML
	doc, err := html.Parse(bytes.NewReader([]byte(htmlData)))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var output []string
	var wordCount int

	// Helper function to traverse the nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// List of tags to include
			switch n.Data {
			case "p", "span", "strong", "em", "b", "i", "u", "mark", "small",
				"sub", "sup", "code", "pre", "blockquote", "q", "abbr", "cite",
				"time", "label", "textarea":
				textContent := getTextContent(n)
				// Add text content only if word count is less than 28000
				words := strings.Fields(textContent)
				if wordCount+len(words) <= 2500 {
					output = append(output, textContent)
					wordCount += len(words)
				} else {
					// Stop appending text once we reach 28000 words
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	// Join and filter out gibberish
	joinedText := strings.Join(output, " ")
	joinedText = filterGibberish(joinedText)

	// If content length is greater than 1000, filter sentences
	if len(joinedText) > 1000 {
		return filterSentencesWithPunctuation(joinedText), nil
	}

	return joinedText, nil
}

// getTextContent retrieves text content from a node and its descendants.
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var textContent []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textContent = append(textContent, getTextContent(c))
	}
	return strings.Join(textContent, " ")
}

// filterGibberish filters out any sequences that don't have at least two words.
func filterGibberish(text string) string {
	// Split the text into words
	words := strings.Fields(text)

	// If there are fewer than two words, return an empty string
	if len(words) < 2 {
		return ""
	}

	// Rebuild the text, ensuring there's at least one space between two words
	var filtered []string
	for i := 0; i < len(words)-1; i++ {
		if len(words[i]) > 0 && len(words[i+1]) > 0 && isValidWord(words[i]) && isValidWord(words[i+1]) {
			filtered = append(filtered, words[i], words[i+1])
			i++ // skip the next word since it's part of the current pair
		}
	}

	return strings.Join(filtered, " ")
}

// isValidWord checks if a word contains at least one letter and is not gibberish.
func isValidWord(word string) bool {
	// Check if the word contains any alphabetic character
	for _, r := range word {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// filterSentencesWithPunctuation filters sentences containing punctuation and with more than 3 words.
func filterSentencesWithPunctuation(text string) string {
	// Split the text into sentences using a simple regex for sentence-ending punctuation
	re := regexp.MustCompile(`[.!?,";]`)
	sentences := strings.Split(text, ".")
	var filteredSentences []string

	for _, sentence := range sentences {
		// Trim spaces and check if the sentence contains punctuation and more than 3 words
		sentence = strings.TrimSpace(sentence)
		words := strings.Fields(sentence)
		if len(words) > 3 && re.MatchString(sentence) {
			filteredSentences = append(filteredSentences, sentence)
		}
	}

	// Join the filtered sentences and return
	return strings.Join(filteredSentences, ". ")
}
