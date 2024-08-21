package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	argsWithProg := os.Args[1:]

	if len(argsWithProg) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(argsWithProg) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := os.Args[1]

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("be sure the second argument is a valid number (max concurrency)")
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("be sure the third argument is a valid number (max pages)")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", baseURL)
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency), // Set a limit of 5 concurrent goroutines
		wg:                 &sync.WaitGroup{},
		maxPages:           int(maxPages),
		visitedPagesCount:  0,
	}
	cfg.wg.Add(1)
	cfg.concurrencyControl <- struct{}{}
	go cfg.crawlPage(baseURL)

	cfg.wg.Wait() // Wait for all goroutines to complete
	printReport(cfg.pages, baseURL)
}
