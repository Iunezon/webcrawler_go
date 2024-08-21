package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("website status not OK")
	}

	header := resp.Header.Get("Content-Type")
	if !strings.Contains(header, "text/html") {
		return "", errors.New("website content is not HTML")
	}

	return string(body), nil
}
