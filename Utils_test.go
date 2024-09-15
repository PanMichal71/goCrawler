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

func TestFixupLink(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		expected string
	}{
		{"No protocol", "example.com", "https://example.com/example.com"},
		{"HTTP", "http://example.com", "https://example.com"},
		{"HTTPS", "https://example.com", "https://example.com"},
		{"Subdomain", "http://sub.example.com", "https://sub.example.com"},
		{"Path", "http://example.com/page", "https://example.com/page"},
		{"Path with query", "http://example.com/page?ref=example.com", "https://example.com/page?ref=example.com"},
		{"Path with fragment", "http://example.com/page#section", "https://example.com/page"},
		{"google link with www", "https://www.google.com", "https://www.google.com"},
		{"Path with fragment", "/page/foo/bar#section", "https://example.com/page/foo/bar"},
		{"Path with fragment", "/page/foo/bar/#section", "https://example.com/page/foo/bar"},
		{"Path with fragment", "page/foo/bar#section", "https://example.com/page/foo/bar"},
	}

	//run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FixupLink("https://example.com", tc.link)
			if result != tc.expected {
				t.Errorf("FixupLink(%q) = %v; want %v", tc.link, result, tc.expected)
			}
		})
	}
}
