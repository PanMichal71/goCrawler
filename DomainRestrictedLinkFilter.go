package main

import (
	"net/url"
	"strings"
)

// DomainRestrictedLinkFilter implements the LinkFilter interface to filter links based on domain restrictions.
type DomainRestrictedLinkFilter struct {
	domain string
}

// NewDomainRestrictedLinkFilter creates a new DomainRestrictedLinkFilter with the specified domain.
func NewDomainRestrictedLinkFilter(domain string) *DomainRestrictedLinkFilter {
	return &DomainRestrictedLinkFilter{domain: domain}
}

// FilterLink checks if the link leads outside the specified domain or to a fragment identifier.
func (d *DomainRestrictedLinkFilter) FilterLink(link string) bool {
	return d.isLinkLeadingOutsideDomain(link) || d.isLinkLeadingToFragmentIdentifier(link)
}

// isLinkLeadingOutsideDomain checks if the link leads outside the specified domain.
func (d *DomainRestrictedLinkFilter) isLinkLeadingOutsideDomain(link string) bool {
	// Parse the link to extract its components.
	parsedURL, err := url.Parse(link)
	if err != nil {
		// If the URL cannot be parsed, it's safer to assume it's external.
		return true
	}

	// If the scheme is "mailto" or "tel", the link points outside the web domain.
	if parsedURL.Scheme == "mailto" || parsedURL.Scheme == "tel" {
		return true
	}

	// Extract the domain from the URL's host.
	host := parsedURL.Host

	// Relative URLs (no host in parsed URL) are considered internal.
	if host == "" && strings.HasPrefix(link, "/") {
		return false
	}

	// Check if the URL's host matches the filter's domain or ends with it
	// to handle subdomains correctly.
	if host == d.domain {
		// The link is internal to the domain.
		return false
	}

	if strings.HasSuffix(host, "www."+d.domain) {
		return false // Internal link
	}

	// By default, consider the link as pointing outside the domain.
	return true
}

// isLinkLeadingToFragmentIdentifier checks if the link contains a fragment identifier (#).
func (d *DomainRestrictedLinkFilter) isLinkLeadingToFragmentIdentifier(link string) bool {
	return strings.Contains(link, "#")
}
