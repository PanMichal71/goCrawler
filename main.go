package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	help := flag.Bool("help", false, "Display this help message")
	outputDir := flag.String("outputDir", ".", "Output directory to store the crawled data")
	ignorePathsArg := flag.String("exlusionPaths", "", "Comma-separated list of paths to ignore")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <url1> <url2> ...\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\nArguments:")
		fmt.Println("  <url1> <url2> ...    List of URLs to crawl")
	}
	flag.Parse()

	if *help || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	urls := flag.Args()
	ignorePaths := strings.Split(*ignorePathsArg, ",")

	var wg sync.WaitGroup
	fileDb := NewFileStorage(*outputDir)
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			fmt.Printf("Crawling: %s\n", url)
			webPage := NewWebPage(&HTTPFetcher{})
			crawler := NewCrawler(webPage, fileDb)
			crawler.Crawl(url, ignorePaths)
		}(url)
	}

	wg.Wait()
	fmt.Println("Completed all crawls.")
}
