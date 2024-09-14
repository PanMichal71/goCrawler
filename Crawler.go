package main

import (
	"fmt"
	"strings"
	"time"
)

type Crawler struct {
	webPage        IWebPage
	contentHandler IContentHandler
	crawledLinks   map[string]bool
	linksToCrawl   []string
	linkFilters    []LinkFilter
	domain         string
}

func NewCrawler(webPage IWebPage, contentHandler IContentHandler) *Crawler {

	return &Crawler{
		webPage:        webPage,
		contentHandler: contentHandler,
		crawledLinks:   make(map[string]bool),
		linksToCrawl:   make([]string, 0),
		linkFilters:    make([]LinkFilter, 0),
	}
}

func (c *Crawler) Crawl(url string, ignorePaths []string) {
	c.domain = url

	c.linkFilters = append(c.linkFilters, NewPathExclusionFilter(ignorePaths))
	c.linkFilters = append(c.linkFilters, NewDomainRestrictedLinkFilter(url))
	c.linkFilters = append(c.linkFilters, &LinkToFileFilter{})

	c.linksToCrawl = append(c.linksToCrawl, url)

	c.crawlImpl(c.domain)
}

func (c *Crawler) crawlImpl(url string) {
	fmt.Printf("Crawling: %s, Links to crawl: %d, Crawled: %d\n", url, len(c.linksToCrawl), len(c.crawledLinks))

	c.linksToCrawl = c.linksToCrawl[:len(c.linksToCrawl)-1]
	c.crawledLinks[url] = true

	htmlContent := c.webPage.Load(url)

	c.contentHandler.HandleContent(url, htmlContent)

	links := c.webPage.GetAllLinks()

	c.processLinks(links)

	if len(c.linksToCrawl) > 0 {
		// TODO: Move this to separate class so I can mock it in tests
		time.Sleep(time.Millisecond * time.Duration(10+time.Now().UnixNano()%1000))
		c.crawlImpl(c.linksToCrawl[len(c.linksToCrawl)-1])
	}
}

func (c *Crawler) fixupDomain(link string) string {
	if strings.HasPrefix(link, "http://") {
		// return link
		return strings.Replace(link, "http://", "https://", 1)
	}

	if strings.HasPrefix(link, "https://") {
		return link
	}

	res := ""
	if strings.HasPrefix(link, "/") {
		res = c.domain + link
	} else {
		res = c.domain + "/" + link
	}

	return res
}

func (c *Crawler) fixupLink(link string) string {
	return strings.TrimSuffix(c.fixupDomain(link), "/")
}

func (c *Crawler) processLinks(links map[string]string) {
	for link := range links {
		shouldCrawl := true
		for _, filter := range c.linkFilters {

			if filter.FilterLink(link) {
				// fmt.Printf("Link %s filtered out by %T\n", link, filter)

				shouldCrawl = false
				break
			}
		}

		if shouldCrawl {
			fixedLink := c.fixupLink(link)

			if _, ok := c.crawledLinks[fixedLink]; !ok {
				found := false
				for _, l := range c.linksToCrawl {
					if l == fixedLink {
						found = true
						break
					}
				}

				if !found {
					fmt.Printf("Adding link to crawl: %s\n", fixedLink)
					c.linksToCrawl = append(c.linksToCrawl, fixedLink)
				}
			}
		}
	}
}
