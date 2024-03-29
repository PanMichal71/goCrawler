package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	// Added for mocking
)

// Fetcher interface defines the behavior for fetching HTML
type Fetcher interface {
	FetchHTML(url string) (string, error)
}

// WebPage implements the logic for working with web pages.
type WebPage struct {
	htmlCache string
	fetcher   Fetcher // Added fetcher dependency
}

// NewWebPage constructor now accepts a fetcher interface
func NewWebPage(fetcher Fetcher) *WebPage {
	return &WebPage{fetcher: fetcher}
}

// Load fetches the HTML content of the given URL and caches it.
func (wp *WebPage) Load(urlToCrawl string) string {
	html, err := wp.fetcher.FetchHTML(urlToCrawl)
	if err != nil {
		// Handle error appropriately; for simplicity, just return an error message.
		return "Failed to load the page: " + err.Error()
	}
	wp.htmlCache = html
	return html
}

// GetAllLinks parses the cached HTML and returns all links as a map.
func (wp *WebPage) GetAllLinks() map[string]string {
	links := make(map[string]string)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wp.htmlCache))
	if err != nil {
		// Handle parsing error; for simplicity, return an empty map.
		return links
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links[href] = s.Text()
		}
	})

	return links
}
