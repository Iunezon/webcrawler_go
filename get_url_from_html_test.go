package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name         string
		htmlBody     string
		rawBaseURL   string
		expectedURLs []string
		expectErr    bool
	}{
		{
			name: "valid HTML with absolute and relative links",
			htmlBody: `<html><body>
						<a href="https://example.com/page1">Page 1</a>
						<a href="/page2">Page 2</a>
					</body></html>`,
			rawBaseURL: "https://example.com",
			expectedURLs: []string{
				"https://example.com/page1",
				"https://example.com/page2",
			},
			expectErr: false,
		},
		{
			name: "HTML with invalid URL",
			htmlBody: `<html><body>
						<a href="::invalid-url">Invalid</a>
						<a href="/valid">Valid</a>
					</body></html>`,
			rawBaseURL: "https://example.com",
			expectedURLs: []string{
				"https://example.com/valid",
			},
			expectErr: false,
		},
		{
			name:         "empty HTML",
			htmlBody:     "",
			rawBaseURL:   "https://example.com",
			expectedURLs: []string{},
			expectErr:    false,
		},
		{
			name:         "invalid base URL",
			htmlBody:     `<html><body><a href="/page">Page</a></body></html>`,
			rawBaseURL:   ":invalid-base-url",
			expectedURLs: nil,
			expectErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualURLs, err := getURLsFromHTML(tc.htmlBody, tc.rawBaseURL)

			// Check if error expectation matches actual error
			if tc.expectErr && err == nil {
				t.Errorf("Test %v - expected an error, but got nil", tc.name)
			} else if !tc.expectErr && err != nil {
				t.Errorf("Test %v - did not expect an error, but got: %v", tc.name, err)
			}

			// Check if the actual URLs match the expected URLs
			if !reflect.DeepEqual(actualURLs, tc.expectedURLs) {
				t.Errorf("Test %v - expected URLs: %v, actual URLs: %v", tc.name, tc.expectedURLs, actualURLs)
			}
		})
	}
}
