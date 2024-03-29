package main

import (
	"testing"
)

func TestDomainRestrictedLinkFilter(t *testing.T) {
	domain := "example.com"
	filter := DomainRestrictedLinkFilter{domain: domain}

	tests := []struct {
		name     string
		link     string
		expected bool // true if the link is expected to be allowed (not filtered out)
	}{
		{"Within domain", "/page", false},
		{"Within domain", "./about", false},                        // TODO: apparently this is not common
		{"Within domain", "http://example.com/about#intro", false}, //TODO: should remove the fragment?

		{"Within domain", "http://example.com/page", false},
		{"Within domain", "http://www.example.com/page", false},
		{"Within domain", "https://example.com/page", false},
		{"Within domain", "#przejdz-do-menu", true},
		{"Subdomain", "http://sub.example.com/page", true},
		{"Outside domain", "http://anotherdomain.com", true},
		{"HTTPS within domain", "https://example.com/secure", false},
		{"Relative path", "/internal/page", false},
		{"Leading to fragment", "http://example.com/page#section", true},
		{"Mailto link", "mailto:user@example.com", true},
		{"Telephone link", "tel:+1234567890", true},
		{"Path outside domain", "http://example.com.redirected.com", true},
		{"Subpath matching domain", "http://anotherdomain.com/page?ref=example.com", true},
		{"Exact match but different TLD", "http://example.co/page", true},
		{"Query parameters", "http://example.com/page?action=view", false},
		{"Path leading to file", "http://example.com/document.pdf", false},
		{"External link containing domain", "http://external.com/?ref=example.com", true},
		{"Subdomain with path", "http://blog.example.com/article", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.FilterLink(tc.link)
			if result != tc.expected {
				t.Errorf("%s: FilterLink(%q) = %v; want %v", tc.name, tc.link, result, tc.expected)
			}
		})
	}
}
