package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"https://centrum.potrafiepomoc.org.pl/",
		"https://potrafiepomoc.org.pl",
	}

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create an instance of the crawler
	// Note: This example assumes you have a constructor for WebPage that matches your requirements

	for _, url := range urls {
		wg.Add(1) // Increment the WaitGroup counter.
		go func(url string) {
			defer wg.Done() // Decrement the counter when the goroutine completes.

			fmt.Printf("Crawling: %s\n", url)
			webPage := NewWebPage(&HTTPFetcher{})
			fileDb := &FileDatabase{}
			crawler := NewCrawler(webPage, fileDb)
			crawler.Crawl(url) // Crawl the URL.
			// Adding sleep to avoid hitting the server too hard, adjust as necessary.
			time.Sleep(time.Millisecond * 500)
		}(url)
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Completed all crawls.")
}
