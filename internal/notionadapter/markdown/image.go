package markdown

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// generateImageMarkup builds the Markdown image syntax.
// If href already contains “data:…”, it’s used verbatim.
func generateImageMarkup(alt, href string) string {
	if strings.HasPrefix(href, "data:") {
		parts := strings.SplitN(href, ",", 2)
		if len(parts) == 2 {
			return fmt.Sprintf("![%s](data:image/png;base64,%s)", alt, parts[1])
		}
	}
	return fmt.Sprintf("![%s](%s)", alt, href)
}

// Image fetches an image (if convertToBase64) and returns a Markdown <img> URL.
func Image(alt, href string, convertToBase64 bool) (string, error) {
	if !convertToBase64 || strings.HasPrefix(href, "data:") {
		return generateImageMarkup(alt, href), nil
	}
	resp, err := http.Get(href)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("![%s](data:image/png;base64,%s)", alt, b64), nil
}

// ImageAsync is like Image but respects context cancellation.
func ImageAsync(ctx context.Context, alt, href string, convertToBase64 bool) (string, error) {
	if !convertToBase64 || strings.HasPrefix(href, "data:") {
		return generateImageMarkup(alt, href), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, href, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("![%s](data:image/png;base64,%s)", alt, b64), nil
}
