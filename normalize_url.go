package main

import (
	"errors"
	"strings"
)

func normalizeURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("no URL provided")
	}

	splittedURL := strings.SplitAfter(url, `:`)

	if len(splittedURL) < 2 {
		return "", errors.New("URL provided is not correct")
	}

	normalizedURL := splittedURL[1]
	countSlash := 0
	for i := 0; i < 2; i++ {
		if string(normalizedURL[i]) == "/" {
			countSlash++
		}
	}

	if !strings.Contains(normalizedURL, ".") {
		return "", errors.New("URL provided is not correct")
	}

	normalizedURL = normalizedURL[countSlash:]

	if strings.LastIndex(normalizedURL, "/") == len(normalizedURL)-1 {
		normalizedURL = normalizedURL[:len(normalizedURL)-1]
	}

	return normalizedURL, nil
}
