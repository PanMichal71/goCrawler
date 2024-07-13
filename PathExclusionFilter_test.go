package main

import (
	"testing"
)

// Test cases for the PathExclusionFilter
func TestPathExclusionFilter(t *testing.T) {
	tests := []struct {
		exclusionPaths []string
		link           string
		expected       bool
	}{
		{[]string{"foo", "admin", "private"}, "http://example.com/foo", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/public", false},
		{[]string{"foo", "admin", "private"}, "http://example.com/admin", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/content", false},
		{[]string{"foo", "admin", "private"}, "http://example.com/content/foo/lorem/ipsum", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/private/secret", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/admin/panel", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/foo/", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/admin/settings", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/private/settings", true},
		{[]string{"foo", "admin", "private"}, "http://example.com/private", true},
		{[]string{"foo/admin", "private"}, "http://example.com/foo/admin/not-private", true},
		{[]string{"foo/admin", "private"}, "http://example.com/foo/admin2/not-private", false},
		{nil, "http://example.com/private", false},

		// Corner cases
		{[]string{}, "http://example.com/any", false},                  // Empty exclusionPaths
		{[]string{"foo"}, "", false},                                   // Empty link
		{[]string{"foo"}, "/", false},                                  // Root path in link
		{[]string{"foo"}, "http://example.com/foo", true},              // Exact match
		{[]string{"foo"}, "http://example.com/Foo", true},              // Case sensitivity
		{[]string{"foo"}, "http://example.com/foo123", false},          // Substring match
		{[]string{"foo"}, "http://example.com/foo/", true},             // Trailing slash in exclusion path
		{[]string{"foo"}, "http://example.com/foo/extra", true},        // Subpath match
		{[]string{"foo"}, "http://example.com/foo/extra", true},        // Trailing slash with subpath match
		{[]string{"admin"}, "http://example.com/administrator", false}, // Substring but not full path
		{[]string{"foo", "admin", "private"}, "http://example.com/fooadminprivate", false},
	}

	for _, test := range tests {
		filter := NewPathExclusionFilter(test.exclusionPaths)
		result := filter.FilterLink(test.link)
		if result != test.expected {
			t.Errorf("FilterLink(%v, %s) = %v; want %v", test.exclusionPaths, test.link, result, test.expected)
		}
	}
}
