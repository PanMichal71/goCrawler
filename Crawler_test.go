package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockIWebPage is a mock type for the IWebPage interface
type MockIWebPage struct {
	mock.Mock
}

// Define methods on the mock that correspond to those of the IWebPage interface
func (m *MockIWebPage) Load(urlToCrawl string) string {
	args := m.Called(urlToCrawl)
	return args.String(0)
}

func (m *MockIWebPage) GetAllLinks() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

// Mock IStorage
type MockIStorage struct {
	mock.Mock
}

func (m *MockIStorage) Save(url string, hash string) error {
	args := m.Called(url, hash)
	return args.Error(0)
}

func (m *MockIStorage) Open(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

func (m *MockIStorage) Close() {
	m.Called()
}

func TestShouldGetHtmlFromGivenUrl(t *testing.T) {
	// Initialize the mock and the object under test
	webPageMock := new(MockIWebPage)
	StorageMock := new(MockIStorage)

	//expect Storage to be opened
	StorageMock.On("Open", "test.db").Return(nil)

	crawler := NewCrawler(webPageMock, StorageMock)

	// Set up expectations
	webPageMock.On("Load", "https://www.google.com").Return("<html><body><a href=\"https://www.google.com\">Google</a></body></html>")

	// Execute the method
	crawler.Crawl("https://www.google.com")

	// Assert that the expectations were met
	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
}
