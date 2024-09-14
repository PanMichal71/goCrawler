package main

import "testing"

func TestNormalizeDomain(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		expected string
	}{
		{"No protocol", "example.com", "example.com"},
		{"HTTP", "http://example.com", "example.com"},
		{"HTTPS", "https://example.com", "example.com"},
		{"Subdomain", "http://sub.example.com", "sub.example.com"},
		{"Path", "http://example.com/page", "example.com"},
		{"Path with query", "http://example.com/page?ref=example.com", "example.com"},
		{"Path with fragment", "http://example.com/page#section", "example.com"},
		{"google link with www", "https://www.google.com", "google.com"},
	}

	//run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NormalizeDomain(tc.link)
			if result != tc.expected {
				t.Errorf("NormalizeDomain(%q) = %v; want %v", tc.link, result, tc.expected)
			}
		})
	}
}

// tests for normalize url
func TestNormalizeUrl(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		expected string
	}{
		{"No protocol", "example.com", "example.com"},
		{"HTTP", "http://example.com", "example.com"},
		{"HTTPS", "https://example.com", "example.com"},
		{"Subdomain", "http://sub.example.com", "sub.example.com"},
		{"Path", "http://example.com/page", "example.com/page"},
		{"Path with query", "http://example.com/page?ref=example.com", "example.com/page"},
		{"Path with fragment", "http://example.com/page/foo/bar#section", "example.com/page/foo/bar"},
		{"Path with fragment", "example.com/page/foo/bar#section", "example.com/page/foo/bar"},
		{"Path with fragment", "www.example.com/page/foo/bar#section", "example.com/page/foo/bar"},
		{"Path with fragment", "https://www.example.com/page/foo/bar#section", "example.com/page/foo/bar"},
		{"Path with fragment", "ftp://www.example.com/page/foo/bar#section", "example.com/page/foo/bar"},
		{"google link with www", "https://www.google.com", "google.com"},
	}

	//run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NormalizeUrl(tc.link)
			if result != tc.expected {
				t.Errorf("NormalizeUrl(%q) = %v; want %v", tc.link, result, tc.expected)
			}
		})
	}
}
