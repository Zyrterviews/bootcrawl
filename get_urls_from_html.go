package main

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// TODO: refactor
func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	var res []string

	// FIXME:
	if htmlBody == "" {
		return nil, nil
		// return nil, errors.New("cannot parse empty string")
	}

	if rawBaseURL == "" {
		return nil, errors.New("base URL cannot be empty")
	}

	reader := strings.NewReader(htmlBody)

	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			hrefAttrIndex := slices.IndexFunc(
				n.Attr,
				func(e html.Attribute) bool {
					return e.Key == "href"
				},
			)

			if hrefAttrIndex == -1 {
				return nil, fmt.Errorf(
					"could not find `href` attribute on <a/>: %s",
					n.Parent.Data,
				)
			}

			link := n.Attr[hrefAttrIndex].Val

			u, err := url.Parse(link)
			if err != nil {
				return nil, err
			}

			if u.Host == "" {
				link, err = url.JoinPath(rawBaseURL, link)
				if err != nil {
					return nil, err
				}
			}

			link, err = url.PathUnescape(link)
			if err != nil {
				return nil, err
			}

			res = append(res, link)
		}
	}

	buf, err := getURLsFromHTML(node.Data, rawBaseURL)

	if buf != nil {
		res = append(res, buf...)
	}

	if err != nil {
		return res, err
	}

	return res, nil
}
