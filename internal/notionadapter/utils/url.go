package utils

import (
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// GetURLTitle attempts to fetch and extract the title from a webpage
// Returns the extracted title or an error if the operation fails
func GetURLTitle(url string) (string, error) {
	// Create a client with a timeout to prevent hanging on slow sites
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make the HTTP request
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response is HTML
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || !isHTML(contentType) {
		return "", nil
	}

	// Extract the title
	return extractTitleFromHTML(resp.Body)
}

// extractTitleFromHTML parses HTML content and extracts the title tag content
func extractTitleFromHTML(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			return n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if title := f(c); title != "" {
				return title
			}
		}
		return ""
	}

	title := f(doc)
	return cleanTitle(title), nil
}

// cleanTitle removes extra whitespace and common suffixes from the title
func cleanTitle(title string) string {
	// Trim whitespace
	title = strings.TrimSpace(title)

	// Remove common suffixes
	suffixes := []string{
		" | Notion",
		" - Notion",
		" – Notion",
	}

	for _, suffix := range suffixes {
		title = strings.TrimSuffix(title, suffix)
	}

	// Replace multiple spaces with a single space
	title = strings.Join(strings.Fields(title), " ")

	return title
}

// isHTML checks if the content type indicates HTML content
func isHTML(contentType string) bool {
	return contentType == "text/html" ||
		contentType == "application/xhtml+xml" ||
		contentType == "application/xml" ||
		contentType == "text/xml" ||
		(len(contentType) >= 9 && contentType[:9] == "text/html")
}
