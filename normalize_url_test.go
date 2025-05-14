package main

import (
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		errorContains string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "handle http scheme",
			inputURL: "http://example.com",
			expected: "example.com",
		},
		{
			name:     "handle no scheme",
			inputURL: "example.com/path",
			expected: "example.com/path",
		},
		{
			name:     "handle query params",
			inputURL: "https://blog.boot.dev/path?a=1&b=2",
			expected: "blog.boot.dev/path?a=1&b=2",
		},
		{
			name:     "handle fragment",
			inputURL: "https://blog.boot.dev/path#section",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "handle empty path",
			inputURL: "https://blog.boot.dev",
			expected: "blog.boot.dev",
		},
		{
			name:     "handle multiple trailing slashes",
			inputURL: "blog.boot.dev/path////",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "handle subdomains",
			inputURL: "https://sub.blog.boot.dev/path",
			expected: "sub.blog.boot.dev/path",
		},
		{
			name:     "handle port",
			inputURL: "http://example.com:8080/path",
			expected: "example.com:8080/path",
		},

		{
			name:     "lowercase capital letters",
			inputURL: "https://BLOG.boot.dev/PATH",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme and capitals and trailing slash",
			inputURL: "http://BLOG.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:          "handle invalid URL",
			inputURL:      `:\\invalidURL`,
			expected:      "",
			errorContains: "could not parse URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				if tc.errorContains != "" &&
					strings.Contains(err.Error(), tc.errorContains) {
					return
				}

				t.Errorf(
					"Test %v - '%s' FAIL: unexpected error: %v",
					i,
					tc.name,
					err,
				)
				return
			}
			if actual != tc.expected {
				t.Errorf(
					"Test %v - %s FAIL: expected URL: %v, actual: %v",
					i,
					tc.name,
					tc.expected,
					actual,
				)
			}
		})
	}
}
