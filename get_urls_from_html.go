package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	htmlParser, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	var URLs []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					hrefURL, err := url.Parse(a.Val)
					if err != nil {
						continue // Skip invalid URLs
					}
					resolvedURL := base.ResolveReference(hrefURL)
					URLs = append(URLs, resolvedURL.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(htmlParser)

	if len(URLs) == 0 {
		return []string{}, nil
	}

	return URLs, nil
}
