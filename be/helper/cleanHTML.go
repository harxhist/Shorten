package helper

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
	"be/logger"
)

var log = logger.Logger;

func CleanContent(content string, tagsToRemove []string) (string, error) {
	doc, err := html.ParseFragment(strings.NewReader(content), nil)
	if err != nil {
		log.Error("Failed to parse HTML content: ", err)
		return "", err
	}

	var buf bytes.Buffer
	// Create a map for quick lookup of tags to remove
	tagsMap := make(map[string]struct{})
	for _, tag := range tagsToRemove {
		tagsMap[tag] = struct{}{}
	}

	// Traverse and clean the HTML tree
	for _, node := range doc {
		cleanNode(node, &buf, tagsMap)
	}

	return buf.String(), nil
}

// cleanNode recursively processes and removes unwanted elements from HTML nodes
func cleanNode(n *html.Node, buf *bytes.Buffer, tagsMap map[string]struct{}) {
	if n.Type == html.ElementNode {
		// Check if the current tag is in the removal list
		if _, remove := tagsMap[n.Data]; remove {
			return
		}
		// Write the opening tag
		buf.WriteString("<" + n.Data)
		// Add any attributes
		for _, attr := range n.Attr {
			buf.WriteString(" " + attr.Key + `="` + attr.Val + `"`)
		}
		buf.WriteString(">")

		// Recursively process child nodes
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			cleanNode(child, buf, tagsMap)
		}

		// Write the closing tag
		buf.WriteString("</" + n.Data + ">")
	} else if n.Type == html.TextNode {
		// Write the text content
		buf.WriteString(n.Data)
	}
}
