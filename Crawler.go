package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// TODO:
// 1. Add exception handling
// 2. Add logging
// 3. Add tests for filters - DONE
// 4. Add Storage to store page hashes - DONE
// 5. Add possibility to interrupt crawling
// 6. Add possibility to configure ignored files
// 7. Add possibility to ignore certain paths, e.g. podopieczni
// 7. Add some generic mechanism to filter links which is configurable from outside
// 8. Add possibility to configure the depth of the crawl
// 9. Add possibility to configure the number of concurrent crawls
// 10. Add possibility to configure the crawl speed
// 11. Add possibility to configure the user agent
// 12. Add possibility to configure the timeout
// 14. Add possibility to configure the maximum number of links to crawl
// 15. Add possibility to read robots.txt and respect it

// I could store md5 from previous run and compare it with the current one
// if they differ I swap them
// At the same time I keep htmls from previous run and save from current one
// so I can compare them and see what has changed
// htmls could be saved in directory tree encompassing the REST paths
// example.com/ -> index.html
// example.com/about -> about.html
// example.com/about/team -> about/team.html

type Crawler struct {
	webPage      IWebPage
	db           IStorage
	crawledLinks map[string]bool // using a map for O(1) access
	linksToCrawl []string
	domain       string
	linkFilters  []LinkFilter
	fileName     string
	file         *os.File
	mutex        sync.Mutex // To safely access linksToCrawl and crawledLinks from multiple goroutines
}

func NewCrawler(webPage IWebPage, db IStorage) *Crawler {

	return &Crawler{
		webPage:      webPage,
		db:           db,
		crawledLinks: make(map[string]bool),
		linksToCrawl: make([]string, 0),
		linkFilters:  make([]LinkFilter, 0),
	}
}

func (c *Crawler) crawlImpl(url string) {
	fmt.Printf("Crawling: %s, Links to crawl: %d, Crawled: %d\n", url, len(c.linksToCrawl), len(c.crawledLinks))

	c.linksToCrawl = c.linksToCrawl[:len(c.linksToCrawl)-1]
	c.crawledLinks[url] = true

	htmlContent := c.webPage.Load(url)
	md5Hash := getMD5Hash(htmlContent)
	c.db.Save(url, md5Hash)

	links := c.webPage.GetAllLinks()

	// for key := range links {
	// 	fmt.Println("\tLink: ", key)
	// }

	c.processLinks(links)

	if len(c.linksToCrawl) > 0 {
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
			c.mutex.Lock()
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

			c.mutex.Unlock()
		}
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (c *Crawler) Crawl(url string) {
	fileName := fmt.Sprintf("md5-%s-%s.json", NormalizeDomain(url), time.Now().Format("2006-01-02:15:04"))
	c.db.Open(fileName)
	defer c.db.Close()

	c.domain = url
	c.linkFilters = append(c.linkFilters, NewDomainRestrictedLinkFilter(url))
	c.linkFilters = append(c.linkFilters, &LinkToFileFilter{})

	fmt.Printf("File: %s\n", fileName)
	c.linksToCrawl = append(c.linksToCrawl, url)

	c.crawlImpl(c.domain)
}
