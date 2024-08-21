package main

import (
	"errors"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "remove scheme",
			inputURL:    "https://blog.boot.dev/path",
			expectedURL: "blog.boot.dev/path",
			expectedErr: nil,
		},
		{
			name:        "remove scheme",
			inputURL:    "https://blog.boot.dev/path/",
			expectedURL: "blog.boot.dev/path",
			expectedErr: nil,
		},
		{
			name:        "remove scheme",
			inputURL:    "http://blog.boot.dev/path/",
			expectedURL: "blog.boot.dev/path",
			expectedErr: nil,
		},
		{
			name:        "remove scheme",
			inputURL:    "http//",
			expectedURL: "",
			expectedErr: errors.New("URL provided is not correct"),
		},
		{
			name:        "remove scheme",
			inputURL:    "",
			expectedURL: "",
			expectedErr: errors.New("no URL provided"),
		},
		{
			name:        "remove scheme",
			inputURL:    "http://",
			expectedURL: "",
			expectedErr: errors.New("URL provided is not correct"),
		},
		{
			name:        "remove scheme",
			inputURL:    "http://fake_url/ciao",
			expectedURL: "",
			expectedErr: errors.New("URL provided is not correct"),
		},
		// add more test cases here
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if actual != tc.expectedURL {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expectedURL, actual)
			}
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Test %v - %s FAIL: expected Error: %v, actual: %v", i, tc.name, tc.expectedErr, err)
			}
		})
	}
}
