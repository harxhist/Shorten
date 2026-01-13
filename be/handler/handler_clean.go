package handler

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"be/logger"
)

var log = logger.Logger

// Mimic a crawler
func fetchContent(url string) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error("failed to create cookie jar: %w", err)
		return "", fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil // Allow redirects
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create request: %w", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-User", "?1")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to execute request: %w", err)
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handling different HTTP status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch content, status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error("failed to create gzip reader: %w", err)
			return "", fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.Close()
	case "deflate":
		reader, err = zlib.NewReader(resp.Body)
		if err != nil {
			log.Error("failed to create zlib reader: %w", err)
			return "", fmt.Errorf("failed to create zlib reader: %w", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	// Reading the response body
	body, err := io.ReadAll(reader)
	if err != nil {
		log.Error("failed to read response body: %w", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}