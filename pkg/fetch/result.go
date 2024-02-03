// Package fetch provides fetching utilities for various content types
/*
 * File: result.go
 * Project: fetch
 * File Created: Sunday, 29th March 2020 5:26:58 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Sunday, 29th March 2020 5:27:21 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package fetch

import (
	"encoding/json"
	"sync"
)

// Result is a struct for holding the results of a fetch action
type Result struct {
	// Query is the original query string
	Query string `json:"query"`
	// Results are the returned formatted Searx results
	Results []SearxResult `json:"results"`
	//ResoltNo is the number of returned results
	ResultNo int `json:"resultno"`
	// PageNo is the current pageNo of the response
	PageNo int `json:"pageno"`
	// Ready indicates if this Result is ready or not
	Ready bool `json:"-"`
	// Any errors associated with this Result
	Errors []error `json:"-"`

	sync.Mutex
}

// HasErrors is a helper function to check if a Result contains errors
func (r *Result) HasErrors() bool {
	if len(r.Errors) > 0 {
		return true
	}
	return false
}

// ToJSON is a function for exporting a Result to JSON
func (r *Result) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}
