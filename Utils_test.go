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
