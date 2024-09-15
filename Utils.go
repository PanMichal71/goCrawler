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

func FixupLink(defaultDomain string, link string) string {
	return strings.TrimSuffix(RemoveFragment(fixupDomain(defaultDomain, link)), "/")
}

func fixupDomain(defaultDomain string, link string) string {
	if strings.HasPrefix(link, "http://") {
		// return link
		return strings.Replace(link, "http://", "https://", 1)
	}

	if strings.HasPrefix(link, "https://") {
		return link
	}

	res := ""
	if strings.HasPrefix(link, "/") {
		res = defaultDomain + link
	} else {
		res = defaultDomain + "/" + link
	}

	return res
}

func RemoveFragment(link string) string {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	parsedUrl.Fragment = ""
	return parsedUrl.String()
}
