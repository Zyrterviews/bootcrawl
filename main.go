package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")

		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")

		os.Exit(1)
	}

	baseURL := os.Args[1]

	var (
		maxPages       int
		maxConcurrency int
	)

	if len(os.Args) > 2 {
		maxPages, _ = strconv.Atoi(os.Args[2])
	}

	if len(os.Args) > 3 {
		maxConcurrency, _ = strconv.Atoi(os.Args[3])
	}

	if maxPages == 0 {
		maxPages = 100
	}

	if maxConcurrency == 0 {
		maxConcurrency = 5
	}

	fmt.Println("starting crawl of:", baseURL)
	fmt.Println("")

	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)

	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            u,
		maxPages:           maxPages,
		mu:                 &mu,
		wg:                 &wg,
		concurrencyControl: make(chan struct{}, maxConcurrency),
	}

	cfg.crawlPage(baseURL)

	cfg.wg.Wait()

	printReport(cfg.pages, baseURL)
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for %s
=============================
        `, baseURL)

	keys := make([]string, 0, len(pages))

	for key := range pages {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		if pages[keys[i]] == pages[keys[j]] {
			return rune(keys[i][0]) > rune(keys[j][0])
		}

		return pages[keys[i]] > pages[keys[j]]
	})

	for _, key := range keys {
		v := pages[key]
		fmt.Printf("Found %d internal links to %s\n", v, key)
	}
}
