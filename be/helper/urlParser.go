package helper

import (
    "strings"
    "errors"
)

// ParseURL normalizes URLs to ensure they start with 'http://' and are correctly formatted
func ParseURL(inputURL string) (string, error) {
    inputURL = strings.TrimSpace(inputURL)
    if inputURL == "" {
		return "", errors.New("URL is empty")
	}

    // Case 1: If it starts with https://www.url.com, convert to http://www.url.com
    if strings.HasPrefix(inputURL, "https://www.") {
        return "http://" + inputURL[len("https://"):], nil;
    }

    // Case 2: If it starts with www.url..., convert to http://www.url.com
    if strings.HasPrefix(inputURL, "www.") {
        return "http://" + inputURL, nil;
    }

    // Case 3: If it starts with https://, convert to http://
    if strings.HasPrefix(inputURL, "https://") {
        return "http://" + inputURL[len("https://"):], nil;
    }

    // Case 4: If it starts with just url.com, convert to http://url.com
    if !strings.HasPrefix(inputURL, "http://") {
        return "http://" + inputURL, nil;
    }

    // Case 5: If it's already in the form of http://www.url.com, maintain it
    return inputURL, nil;
}