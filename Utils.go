package main

import (
	"fmt"
	"net/url"
)

func normalizeDomain(domain string) string {
	// Parse the URL string
	parsedUrl, err := url.Parse(domain)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return domain // Handle invalid URLs gracefully
	}

	// if hostname is empty, return the domain as is else return hostname
	if parsedUrl.Hostname() == "" {
		return domain
	}

	return parsedUrl.Hostname()
}
