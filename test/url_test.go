package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kible/notion-to-md/internal/notionadapter/utils"
)

func TestGetURLTitle(t *testing.T) {
	// Create a test server that returns a simple HTML page with a title
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Test Page Title</title>
			</head>
			<body>
				<h1>Hello, world!</h1>
			</body>
			</html>
		`))
	}))
	defer ts.Close()

	// Test successful title extraction
	title, err := utils.GetURLTitle(ts.URL)
	if err != nil {
		t.Errorf("GetURLTitle returned an error: %v", err)
	}
	if title != "Test Page Title" {
		t.Errorf("GetURLTitle returned wrong title, got: %s, want: %s", title, "Test Page Title")
	}

	// Create a test server that returns non-HTML content
	tsNonHTML := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "This is not HTML"}`))
	}))
	defer tsNonHTML.Close()

	// Test non-HTML content
	title, err = utils.GetURLTitle(tsNonHTML.URL)
	if err != nil {
		t.Errorf("GetURLTitle returned an error for non-HTML content: %v", err)
	}
	if title != "" {
		t.Errorf("GetURLTitle should return empty string for non-HTML content, got: %s", title)
	}

	// Test invalid URL
	title, err = utils.GetURLTitle("http://invalid-url-that-does-not-exist.example")
	if err == nil {
		t.Error("GetURLTitle should return an error for invalid URL")
	}
}
