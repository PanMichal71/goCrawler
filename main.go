package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	// Define a help flag
	help := flag.Bool("help", false, "Display this help message")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <url1> <url2> ...\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\nArguments:")
		fmt.Println("  <url1> <url2> ...    List of URLs to crawl")
	}
	flag.Parse()

	// Display help message if the help flag is set or no URLs are provided
	if *help || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	urls := flag.Args()

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1) // Increment the WaitGroup counter.
		go func(url string) {
			defer wg.Done() // Decrement the counter when the goroutine completes.

			fmt.Printf("Crawling: %s\n", url)
			webPage := NewWebPage(&HTTPFetcher{})
			fileDb := &FileStorage{}
			crawler := NewCrawler(webPage, fileDb)
			crawler.Crawl(url) // Crawl the URL.
		}(url)
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Completed all crawls.")
}
