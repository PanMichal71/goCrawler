package main

// IWebPage defines the interface for web page operations.
type IWebPage interface {
    Load(urlToCrawl string) string
    GetAllLinks() map[string]string
}
