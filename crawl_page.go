package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
	visitedPagesCount  int
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()      // Ensure that the WaitGroup counter is decremented when the goroutine finishes
	<-cfg.concurrencyControl // Release the spot in the concurrency channel when the function is done

	if cfg.visitedPagesCount >= cfg.maxPages {
		return
	}

	currentDomain, err := ExtractDomain(rawCurrentURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Ensure that the current domain matches the base domain
	baseDomain := cfg.baseURL.Hostname()
	if baseDomain != currentDomain {
		return
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	isFirst := cfg.addPageVisit(normCurrentURL)
	if !isFirst {
		return // If we've already visited this page, skip it
	}

	websiteBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(rawCurrentURL, ":", err.Error())
		return
	}

	URLs, err := getURLsFromHTML(websiteBody, cfg.baseURL.String())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, url := range URLs {
		var prefix string
		if strings.HasPrefix(url, "http://") {
			prefix = "http://"
		} else if strings.HasPrefix(url, "https://") {
			prefix = "https://"
		}

		normalizedURL, err := normalizeURL(url)
		if err != nil {
			continue
		}

		cfg.mu.Lock()
		cfg.wg.Add(1)                            // Add a goroutine to the WaitGroup
		cfg.concurrencyControl <- struct{}{}     // Block if the channel is full
		go cfg.crawlPage(prefix + normalizedURL) // Start a new goroutine to crawl the next page
		cfg.mu.Unlock()
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normalizedURL]; ok {
		cfg.pages[normalizedURL] += 1
		return false
	}

	cfg.pages[normalizedURL] = 1
	cfg.visitedPagesCount++
	return true
}

func ExtractDomain(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	host := parsedURL.Host

	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	return host, nil
}

/*
	func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
		baseDomain, err := ExtractDomain(rawBaseURL)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		currentDomain, err := ExtractDomain(rawCurrentURL)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if baseDomain != currentDomain {
			return
		}

		normCurrentURL, err := normalizeURL(rawCurrentURL)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, ok := pages[normCurrentURL]
		if ok {
			pages[normCurrentURL] += 1
		} else {
			pages[normCurrentURL] = 1
		}

		websiteBody, err := getHTML(rawCurrentURL)
		if err != nil {
			fmt.Println(rawCurrentURL, ":", err.Error())
			return
		}

		URLs, err := getURLsFromHTML(websiteBody, rawBaseURL)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, url := range URLs {
			var prefix string
			if strings.HasPrefix(url, "http://") {
				prefix = "http://"
			} else if strings.HasPrefix(url, "https://") {
				prefix = "https://"
			}
			url, _ = normalizeURL(url)
			if _, ok := pages[url]; !ok {
				crawlPage(rawBaseURL, prefix+url, pages)
			} else {
				pages[url] += 1
			}
		}

}
*/
