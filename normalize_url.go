package main

import (
	"fmt"
	"net/url"
	"strings"
)

// normalizeURL removes the scheme (http:// or https://), trailing slashes,
// and fragments from a URL, returning the normalized string
func normalizeURL(inputURL string) (string, error) {
	inputURL = strings.ToLower(inputURL)

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("could not parse URL: %w", err)
	}

	host := parsedURL.Host

	if host == "" {
		//nolint:mnd
		parts := strings.SplitN(parsedURL.Path, "/", 2)
		host = parts[0]

		if len(parts) > 1 {
			parsedURL.Path = "/" + parts[1]
		} else {
			parsedURL.Path = ""
		}
	}

	path := strings.TrimRight(parsedURL.Path, "/")

	var query string

	if parsedURL.RawQuery != "" {
		query = "?" + parsedURL.RawQuery
	}

	return host + path + query, nil
}
