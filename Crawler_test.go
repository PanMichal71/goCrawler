package main

import (
	"regexp"
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

func fileNameMatchesPattern() interface{} {
	return func(fileName string) bool {
		// Define the regex pattern to match the filename format.
		// regexPattern := `^md5-.+-\d{4}-\d{2}-\d{2}:\d{2}:\d{2}(\-.+)?\.json$`
		regexPattern := `^md5-.+-\d{4}-\d{2}-\d{2}:\d{2}:\d{2}:\d{2}(\-.+)?\.json$`

		matched, err := regexp.MatchString(regexPattern, fileName)
		if err != nil {
			return false
		}
		return matched
	}
}

// define variable for all the tests to store this string: "<html><body><a href=\"https://www.google.com\">Google</a></body></html>"
var defaultHtmlContent = "<html><body><a href=\"https://www.google.com\">Google</a></body></html>"
var defaultHtmlContentMd5Hash = "d6165a2f6a47eba8aa611ca6891203a9"

func TestShouldGetHtmlFromGivenUrl(t *testing.T) {
	// Initialize the mock and the object under test
	webPageMock := new(MockIWebPage)
	storageMock := new(MockIStorage)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern())).Return(nil)
	storageMock.On("Save", "https://www.google.com", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Close").Return()

	crawler := NewCrawler(webPageMock, storageMock)

	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	//expect call for GetAllLinks and it returns empty map
	webPageMock.On("GetAllLinks").Return(map[string]string{})

	// Execute the method
	crawler.Crawl("https://www.google.com", nil)

	// Assert that the expectations were met
	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
}

func TestShouldOnlyCrawlNotVisitedLinks(t *testing.T) {
	webPageMock := new(MockIWebPage)
	storageMock := new(MockIStorage)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern())).Return(nil)
	storageMock.On("Save", "https://www.google.com", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Save", "https://www.google.com/kontakty", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Close").Return()

	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("Load", "https://www.google.com/kontakty").Return(defaultHtmlContent)

	webPageMock.On("GetAllLinks").Return(map[string]string{
		"https://www.google.com": "l1",
		"/kontakty":              "l3",
	}).Once()
	webPageMock.On("GetAllLinks").Return(map[string]string{
		"l1": "https://www.google.com"}).Once()
	webPageMock.On("GetAllLinks").Return(map[string]string{}).Maybe()

	crawler := NewCrawler(webPageMock, storageMock)
	crawler.Crawl("https://www.google.com", nil)

	webPageMock.AssertNumberOfCalls(t, "GetAllLinks", 2)
	webPageMock.AssertCalled(t, "Load", "https://www.google.com/kontakty")
}

func TestShouldOnlyCrawlToSameDomain(t *testing.T) {
	webPageMock := new(MockIWebPage)
	storageMock := new(MockIStorage)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern())).Return(nil)
	storageMock.On("Save", "https://www.google.com", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Save", "https://www.google.com/kontakty", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Close").Return()

	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("Load", "https://www.google.com/kontakty").Return(defaultHtmlContent)

	webPageMock.On("GetAllLinks").Return(map[string]string{
		"https://www.google.com/kontakty": "l2",
		"/kontakty":                       "l3",
	}).Once()
	webPageMock.On("GetAllLinks").Return(map[string]string{}).Maybe()

	crawler := NewCrawler(webPageMock, storageMock)
	crawler.Crawl("https://www.google.com", nil)

	webPageMock.AssertCalled(t, "Load", "https://www.google.com/kontakty")
	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
	webPageMock.AssertNumberOfCalls(t, "Load", 2)
}

func TestShouldFilterOutLinksLeadingOutsideOfDomain(t *testing.T) {
	webPageMock := new(MockIWebPage)
	storageMock := new(MockIStorage)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern())).Return(nil)
	storageMock.On("Save", "https://www.google.com", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Close").Return()

	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("GetAllLinks").Return(map[string]string{
		"tel:+48509685328":                  "l1",
		"https://www.google2.com":           "l2",
		"mailto:sekretariat@example.org.pl": "l3",
	}).Once()

	crawler := NewCrawler(webPageMock, storageMock)
	crawler.Crawl("https://www.google.com", nil)

	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
	webPageMock.AssertNumberOfCalls(t, "Load", 1)
}

func TestShouldNotAddLinksAlreadyInQueue(t *testing.T) {
	webPageMock := new(MockIWebPage)
	storageMock := new(MockIStorage)

	storageMock.On("Open", mock.MatchedBy(fileNameMatchesPattern())).Return(nil)
	storageMock.On("Save", "https://www.google.com", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Save", "https://www.google.com/pomoc", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Save", "https://www.google.com/kontakty", defaultHtmlContentMd5Hash).Return(nil)
	storageMock.On("Close").Return()

	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("Load", "https://www.google.com/pomoc").Return(defaultHtmlContent)
	webPageMock.On("Load", "https://www.google.com/kontakty").Return(defaultHtmlContent)
	webPageMock.On("GetAllLinks").Return(map[string]string{
		"https://www.google.com":        "l1",
		"https://www.google.com/pomoc":  "l2",
		"/kontakty":                     "l3",
		"http://www.google.com":         "l4",
		"www.google.com":                "l5",
		"google.com":                    "l6",
		"http://www.google.com/":        "l7",
		"https://www.google.com/pomoc/": "l8",
	}).Once()

	webPageMock.On("GetAllLinks").Return(map[string]string{
		"/kontakty": "l1",
	}).Once()
	webPageMock.On("GetAllLinks").Return(map[string]string{}).Once()

	crawler := NewCrawler(webPageMock, storageMock)
	crawler.Crawl("https://www.google.com", nil)

	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
	webPageMock.AssertCalled(t, "Load", "https://www.google.com/pomoc")
	webPageMock.AssertCalled(t, "Load", "https://www.google.com/kontakty")
	webPageMock.AssertNumberOfCalls(t, "Load", 3)
}
