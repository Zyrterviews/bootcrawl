package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("URL cannot be empty")
	}

	res, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf(
			"an error occurred while fetching the webpage: %s",
			res.Status,
		)
	}

	contentType := strings.Split(res.Header.Get("Content-Type"), ";")
	// TODO: handle `charset=xxx`, etc.. maybe a `strings.Contains` but are you sure?
	if contentType[0] != "text/html" {
		return "", fmt.Errorf(
			"expected Content-Type to be `text/html`, got `%s` instead",
			contentType[0],
		)
	}

	htmlBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if len(htmlBody) == 0 {
		return "", errors.New("received empty HTML body")
	}

	return string(htmlBody), nil
}
