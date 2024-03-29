package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPFetcher implements the Fetcher interface using HTTP GET requests.
type HTTPFetcher struct {
	Client *http.Client // Optional: Allow customization of the HTTP client
}

// NewHTTPFetcher creates a new HTTPFetcher instance with a default HTTP client.
func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{Client: http.DefaultClient}
}

// SetClient allows customization of the HTTP client used by the fetcher.
func (f *HTTPFetcher) SetClient(client *http.Client) {
	f.Client = client
}

// FetchHTML makes an HTTP GET request to the specified URL and returns the body as a string.
func (f *HTTPFetcher) FetchHTML(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
