package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	maxPages           int
	concurrencyControl chan struct{}
	mu                 *sync.Mutex
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	if cfg.baseURL == nil {
		return
	}

	if rawCurrentURL == "" {
		return
	}

	cfg.mu.Lock()

	if cfg.pages == nil || len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()

		return
	}

	cfg.mu.Unlock()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)

		return
	}

	if cfg.baseURL.Host != currentURL.Host {
		fmt.Println("not same host")

		return
	}

	nu, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)

		return
	}

	cfg.mu.Lock()

	if _, ok := cfg.pages[nu]; ok {
		fmt.Println("Already got:", nu)

		cfg.pages[nu]++

		cfg.mu.Unlock()

		return
	}

	cfg.pages[nu] = 1

	cfg.mu.Unlock()

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(htmlBody)

	urls, err := getURLsFromHTML(htmlBody, cfg.baseURL.String())
	if err != nil {
		fmt.Println(err)

		return
	}

	for _, u := range urls {
		cfg.wg.Add(1)

		cfg.concurrencyControl <- struct{}{}

		go func() {
			defer cfg.wg.Done()

			<-cfg.concurrencyControl
			cfg.crawlPage(u)
		}()
	}
}
