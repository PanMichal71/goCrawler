package main

import (
	"testing"
)

// TestLinkToFileFilter tests the FilterLink method of LinkToFileFilter.
func TestLinkToFileFilter(t *testing.T) {
	// Create an instance of LinkToFileFilter.
	filter := LinkToFileFilter{}

	// Define test cases
	tests := []struct {
		name     string // Added test case name for better identification
		link     string
		expected bool
	}{
		{"PDF extension", "http://example.com/file.pdf", true},
		{"PNG extension", "http://example.com/image.png", true},
		{"No extension", "http://example.com", false},
		{"Text after dot, but not extension", "http://example.com/doc", false},
		{"JPG extension", "http://example.com/photo.jpg", true},
		{"Uppercase extension", "http://example.com/image.JPG", true},                      // Test case for uppercase extension
		{"Query parameters", "http://example.com/download?file=report.pdf", false},         // Depending on implementation
		{"Fragment identifier", "http://example.com/document.pdf#chapter1", true},          // Depending on implementation
		{"Mimic file extension in directory", "http://example.com/.pdfs/info", false},      // Directory that mimics extension
		{"Path and query", "http://example.com/files/download.jpg?action=download", true},  // Path with extension and query
		{"Path and query", "http://example.com/files/download.JPG?action=download", true},  // Path with extension in upper case and query
		{"Nested path with extension", "http://example.com/archive/2021/report.pdf", true}, // Nested path with extension
		{"URL with fragment", "http://example.com/files/report.pdf#page=2", true},
		{"URL-encoded path", "http://example.com/files/document%20proposal.pdf", true},
		{"Complex path with multiple extensions", "http://example.com/archive.2021/images.tar.gz", true},
		{"Non-file URLs", "http://example.com/login", false},
		{"Query parameters without values", "http://example.com/download?file", false},
		{"Subdomains and ports", "http://sub.example.com/files/info.pdf", true},
		{"URLs ending with a slash", "http://example.com/documents/", false},
		{"Multiple file extensions", "http://example.com/files/compressed.tar.gz", true},
		{"Multiple file extensions", "http://example.com/files/compressed.tar.gz.zip", true},
		{"Multiple file extensions", "http://example.com/files/compressed.tar.gz.zip?action=download", true},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { // Use t.Run for subtests
			result := filter.FilterLink(tc.link)
			if result != tc.expected {
				t.Errorf("FilterLink(%q) = %v; want %v", tc.link, result, tc.expected)
			}
		})
	}
}
