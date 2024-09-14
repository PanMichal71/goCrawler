package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockIWebPage struct {
	mock.Mock
}

func (m *MockIWebPage) Load(urlToCrawl string) string {
	args := m.Called(urlToCrawl)
	return args.String(0)
}

func (m *MockIWebPage) GetAllLinks() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

type MockIContentHandler struct {
	mock.Mock
}

func (m *MockIContentHandler) HandleContent(url string, content string) error {
	m.Called(url, content)
	return nil
}

func TestShouldGetHtmlFromGivenUrl(t *testing.T) {
	contentHandlerMock := new(MockIContentHandler)
	contentHandlerMock.On("HandleContent", "https://www.google.com", defaultHtmlContent).Return()

	webPageMock := new(MockIWebPage)
	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("GetAllLinks").Return(map[string]string{})

	crawler := NewCrawler(webPageMock, contentHandlerMock)
	crawler.Crawl("https://www.google.com", nil)

	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
}

func TestShouldOnlyCrawlNotVisitedLinks(t *testing.T) {
	contentHandlerMock := new(MockIContentHandler)
	contentHandlerMock.On("HandleContent", "https://www.google.com", defaultHtmlContent).Return()
	contentHandlerMock.On("HandleContent", "https://www.google.com/kontakty", defaultHtmlContent).Return()

	webPageMock := new(MockIWebPage)
	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("Load", "https://www.google.com/kontakty").Return(defaultHtmlContent)

	webPageMock.On("GetAllLinks").Return(map[string]string{
		"https://www.google.com": "l1",
		"/kontakty":              "l3",
	}).Once()
	webPageMock.On("GetAllLinks").Return(map[string]string{
		"l1": "https://www.google.com"}).Once()
	webPageMock.On("GetAllLinks").Return(map[string]string{}).Maybe()

	crawler := NewCrawler(webPageMock, contentHandlerMock)
	crawler.Crawl("https://www.google.com", nil)

	webPageMock.AssertNumberOfCalls(t, "GetAllLinks", 2)
	webPageMock.AssertCalled(t, "Load", "https://www.google.com/kontakty")
}

func TestShouldOnlyCrawlToSameDomain(t *testing.T) {
	contentHandlerMock := new(MockIContentHandler)
	contentHandlerMock.On("HandleContent", "https://www.google.com", defaultHtmlContent).Return()
	contentHandlerMock.On("HandleContent", "https://www.google.com/kontakty", defaultHtmlContent).Return()

	webPageMock := new(MockIWebPage)
	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("Load", "https://www.google.com/kontakty").Return(defaultHtmlContent)

	webPageMock.On("GetAllLinks").Return(map[string]string{
		"https://www.google.com/kontakty": "l2",
		"/kontakty":                       "l3",
	}).Once()
	webPageMock.On("GetAllLinks").Return(map[string]string{}).Maybe()

	crawler := NewCrawler(webPageMock, contentHandlerMock)
	crawler.Crawl("https://www.google.com", nil)

	webPageMock.AssertCalled(t, "Load", "https://www.google.com/kontakty")
	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
	webPageMock.AssertNumberOfCalls(t, "Load", 2)
}

func TestShouldFilterOutLinksLeadingOutsideOfDomain(t *testing.T) {
	contentHandlerMock := new(MockIContentHandler)
	contentHandlerMock.On("HandleContent", "https://www.google.com", defaultHtmlContent).Return()

	webPageMock := new(MockIWebPage)
	webPageMock.On("Load", "https://www.google.com").Return(defaultHtmlContent)
	webPageMock.On("GetAllLinks").Return(map[string]string{
		"tel:+48509685328":                  "l1",
		"https://www.google2.com":           "l2",
		"mailto:sekretariat@example.org.pl": "l3",
	}).Once()

	crawler := NewCrawler(webPageMock, contentHandlerMock)
	crawler.Crawl("https://www.google.com", nil)

	contentHandlerMock.AssertNotCalled(t, "HandleContent", "tel:+48509685328", defaultHtmlContent)
	contentHandlerMock.AssertNotCalled(t, "HandleContent", "https://www.google2.com", defaultHtmlContent)
	contentHandlerMock.AssertCalled(t, "HandleContent", "https://www.google.com", defaultHtmlContent)

	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
	webPageMock.AssertNumberOfCalls(t, "Load", 1)
}

func TestShouldNotAddLinksAlreadyInQueue(t *testing.T) {
	contentHandlerMock := new(MockIContentHandler)
	contentHandlerMock.On("HandleContent", "https://www.google.com", defaultHtmlContent).Return()
	contentHandlerMock.On("HandleContent", "https://www.google.com/pomoc", defaultHtmlContent).Return()
	contentHandlerMock.On("HandleContent", "https://www.google.com/kontakty", defaultHtmlContent).Return()

	webPageMock := new(MockIWebPage)
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

	crawler := NewCrawler(webPageMock, contentHandlerMock)
	crawler.Crawl("https://www.google.com", nil)

	contentHandlerMock.AssertCalled(t, "HandleContent", "https://www.google.com", defaultHtmlContent)
	contentHandlerMock.AssertCalled(t, "HandleContent", "https://www.google.com/pomoc", defaultHtmlContent)
	contentHandlerMock.AssertCalled(t, "HandleContent", "https://www.google.com/kontakty", defaultHtmlContent)
	contentHandlerMock.AssertNotCalled(t, "HandleContent", "http://www.google.com", defaultHtmlContent)

	webPageMock.AssertCalled(t, "Load", "https://www.google.com")
	webPageMock.AssertCalled(t, "Load", "https://www.google.com/pomoc")
	webPageMock.AssertCalled(t, "Load", "https://www.google.com/kontakty")
	webPageMock.AssertNumberOfCalls(t, "Load", 3)
}
