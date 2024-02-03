// Package fetch provides fetching utilities for various content types
/*
 * File: fetch.go
 * Project: fetch
 * File Created: Wednesday, 18th March 2020 8:37:31 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Thursday, 28th May 2020 10:20:46 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package fetch

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Fetcher is a struct for holding a fetcher's context variables
// We only allow one type per fetcher instance here so that each
// Fetcher can be configured individually for a specific content type.
type Fetcher struct {
	// Searx connection address
	SearxAddr string
	// Type of search to execute
	Type string
	// Languages to search in
	Languages []string
	// Stream results as they are received
	StreamResults bool

	// HashMap used for filtering out duplicate image sources
	cache map[string]bool
	sync.Mutex

	// TODO: add in image fuzzy de-duping: https://godoc.org/github.com/rivo/duplo#Match
}

// NewFetcher creates a new Fetcher object
func NewFetcher(
	searxAddr, contentType string,
	languagesAll, languagesSimplified, streamResults bool,
	languages ...string) (*Fetcher, error) {

	// HealthCheck
	if !searxHealthCheck(searxAddr) {
		return nil, fmt.Errorf("unable to reach Searx at '%s'. Is a Searx server running at the specified address?", searxAddr)
	}

	// Check content type
	contentType = strings.ToLower(contentType)
	if !CheckIsAvailableContentType(contentType) {
		return nil, fmt.Errorf("unrecognizable content type: '%s'", contentType)
	}

	// Configure language codes
	// Language preferences override in the order: languagesSimplified < languagesAll < explicitly set languages
	var langCodes []string

	if languagesSimplified {
		langCodes = getSimplifiedLanguageCodes()
	}

	if languagesAll {
		langCodes = getAllLanguageCodes()
	}

	if len(languages) > 0 {
		langCodes = nil
		for _, lang := range languages {
			if code, ok := getLanguageCode(lang); ok {
				langCodes = append(langCodes, code)
			} else {
				log.Warnf("unrecognized language definition '%s'; skipping\n", lang)
			}
		}
	}

	if len(langCodes) == 0 {
		return nil, fmt.Errorf("unable to initialize Fetcher; no configured language definitions")
	}

	return &Fetcher{
		SearxAddr:     searxAddr,
		Type:          contentType,
		Languages:     langCodes,
		StreamResults: streamResults,
		cache:         make(map[string]bool),
	}, nil
}

// FetchAsync is a method for asynchronously executing a content query
// for the first or a specified set of results.
func (f *Fetcher) FetchAsync(query string, pageNo ...int) *Result {
	if len(pageNo) == 0 {
		pageNo[0] = 1
	} else {
		if pageNo[0] == 0 {
			pageNo[0] = 1
		}
	}

	result := &Result{
		Query:    query,
		ResultNo: 0,
		PageNo:   pageNo[0],
		Ready:    false,
	}

	var errs []error
	var wg sync.WaitGroup

	// Queue up a goroutine for each lang code
	for _, lang := range f.Languages {

		wg.Add(1)
		go func(lang string, r *Result, wg *sync.WaitGroup) {
			defer wg.Done()

			log.Infof("search: query=%s, type=%s, lang=%s, pageNo=%d\n", query, f.Type, lang, pageNo[0])

			results, err := searxQuery(f.SearxAddr, query, f.Type, lang, pageNo[0])
			if err != nil {
				errs = append(errs, fmt.Errorf("unexpected error in Searx query with query '%s' (lang=%s); err=%s", query, lang, err))
			} else {
				// Filter results using url cache
				filteredResults := []SearxResult{}

				f.Lock()
				for _, r := range results {
					if _, ok := f.cache[r.URL]; !ok {
						f.cache[r.URL] = true
						filteredResults = append(filteredResults, r)
					} else {
						log.Debugf("filtering out: %s", r.URL)
					}
				}
				f.Unlock()

				if len(filteredResults) == 0 {
					return
				}

				// Update result
				result.Lock()
				result.Results = append(result.Results, filteredResults...)
				result.ResultNo += len(filteredResults)
				result.Unlock()

				// Stream if specified
				if f.StreamResults {
					for _, res := range filteredResults {
						j, err := json.MarshalIndent(res, "", "\t")
						if err != nil {
							continue
						}
						fmt.Fprint(os.Stdout, string(j)+"\n")
					}
				}
			}
		}(lang, result, &wg)
	}

	// Goroutine for updating fetch.Result Ready & Errors attributes
	go func(r *Result, wg *sync.WaitGroup) {
		wg.Wait()
		result.Lock()
		result.Ready = true
		result.Errors = append(result.Errors, errs...)
		result.Unlock()

	}(result, &wg)

	return result
}

// ClearCache is a function to clear a Fetcher's URL cache
func (f *Fetcher) ClearCache() {
	for k := range f.cache {
		delete(f.cache, k)
	}
}

// CacheLen is a function to retrieve a Fetcher's URL cache length
func (f *Fetcher) CacheLen() int {
	return len(f.cache)
}
