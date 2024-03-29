package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockWebPage struct {
	WebPage  // Embedding WebPage to inherit its methods.
	mockHTML string
}

type mockFetcher struct {
	mock.Mock
}

func (m *mockFetcher) FetchHTML(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func TestGetAllLinks(t *testing.T) {
	mockFetcher := new(mockFetcher)
	mockFetcher.On("FetchHTML", "https://example.com").Return(`<html><body><a href="http://example.com">Example</a><a href="http://test.com">Test</a></body></html>`, nil)

	wp := NewWebPage(mockFetcher)

	wp.Load("https://example.com")
	links := wp.GetAllLinks()

	expectedLinks := map[string]string{
		"http://example.com": "Example",
		"http://test.com":    "Test",
	}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("Expected links to be %v, got %v instead", expectedLinks, links)
	}
}

func TestGetAllLinks_EmptyHTML(t *testing.T) {
	mockFetcher := new(mockFetcher)
	mockFetcher.On("FetchHTML", "https://example.com").Return("", nil)

	wp := NewWebPage(mockFetcher)
	wp.Load("https://example.com")
	links := wp.GetAllLinks()

	expectedLinks := map[string]string{}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("Expected links to be empty, got %v instead", links)
	}
}

func TestGetAllLinks_ParsingError(t *testing.T) {
	mockFetcher := new(mockFetcher)
	mockFetcher.On("FetchHTML", "https://example.com").Return("", fmt.Errorf("parsing error"))

	wp := NewWebPage(mockFetcher)
	wp.Load("https://example.com")
	links := wp.GetAllLinks()

	expectedLinks := map[string]string{}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("Expected links to be empty, got %v instead", links)
	}
}

func TestGetAllLinks_MissingLinks(t *testing.T) {
	mockFetcher := new(mockFetcher)
	mockFetcher.On("FetchHTML", "https://example.com").Return(`<html><body></body></html>`, nil)

	wp := NewWebPage(mockFetcher)
	wp.Load("https://example.com")
	links := wp.GetAllLinks()

	expectedLinks := map[string]string{}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("Expected links to be empty, got %v instead", links)
	}
}

func TestGetAllLinks_RelativeLinks(t *testing.T) {
	mockFetcher := new(mockFetcher)
	mockFetcher.On("FetchHTML", "https://example.com/base").Return(`<html><body><a href="./about">About</a></body></html>`, nil)

	wp := NewWebPage(mockFetcher)
	wp.Load("https://example.com/base")
	links := wp.GetAllLinks()

	expectedLinks := map[string]string{
		"./about": "About",
	}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("Expected links to be %v, got %v instead", expectedLinks, links)
	}
}

func TestGetAllLinks_FragmentLinks(t *testing.T) {
	// Setup: Mock with HTML containing fragment links (optional test)
	mockFetcher := new(mockFetcher)
	mockFetcher.On("FetchHTML", "https://example.com").Return(`<html><body><a href="http://example.com/about#intro">Intro</a><a href="https://test.com/contact">Contact</a></body></html>`, nil)

	wp := NewWebPage(mockFetcher)
	wp.Load("https://example.com")
	links := wp.GetAllLinks()

	expectedLinks := map[string]string{
		"http://example.com/about#intro": "Intro",
		"https://test.com/contact":       "Contact",
	}

	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("Expected links to be %v, got %v instead", expectedLinks, links)
	}
}
