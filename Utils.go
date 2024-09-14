package main

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeDomain(domain string) string {
	parsedUrl, err := url.Parse(domain)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	// if hostname is empty, return the domain as is else return hostname
	if parsedUrl.Hostname() == "" {
		return domain
	}

	return strings.TrimPrefix(parsedUrl.Hostname(), "www.")
}

func NormalizeUrl(urlString string) string {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	if parsedUrl.Hostname() == "" {
		return strings.TrimPrefix(parsedUrl.Path, "www.")
	}

	return NormalizeDomain(urlString) + parsedUrl.Path
}
