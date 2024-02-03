// Package fetch provides fetching utilities for various content types
/*
 * File: searx.go
 * Project: fetch
 * File Created: Saturday, 21st March 2020 9:01:55 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Thursday, 28th May 2020 9:52:07 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package fetch

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// searxHealthCheck is a helper function to check connection status to the Searx server
func searxHealthCheck(searxAddr string) bool {
	resp, err := http.Get(searxAddr)
	if err == nil {
		if resp.StatusCode == http.StatusOK {
			return true
		}
	}

	return false
}

// searxQuery is a function for executing a Searx query
func searxQuery(addr string, query string, category string, lang string, pageno int) ([]SearxResult, error) {
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("categories", category)
	q.Add("language", lang)
	q.Add("format", "json")
	q.Add("pageno", fmt.Sprintf("%d", pageno))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	rawResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return nil, err
	}

	fmtResponse, err := formatSearxResponse(body)
	if err != nil {
		return nil, err
	}

	return fmtResponse.Results, nil
}

// SearxResult is a struct for representing the details of a Searx result
type SearxResult struct {
	URL             string `json:"url"`
	ImgFmt          string `json:"img_format"`
	ImgSrc          string `json:"img_src"`
	ImgSrcB64       string `json:"img_src_b64"`
	ThumbnailSrc    string `json:"thumbnail_src"`
	ThumbnailSrcB64 string `json:"thumbnail_src_b64"`
	Engine          string `json:"engine"`
	Source          string `json:"source"`
	Title           string `json:"title"`
}

// searxResponse is a struct for representing a Searx response
type searxResponse struct {
	Results []SearxResult `json:"results"`
}

// formatSearxResponse is a function for unmarshalling a Searx query into a response object
func formatSearxResponse(searxJSONResponse []byte) (*searxResponse, error) {
	var resp searxResponse
	var results = []SearxResult{}

	// Unmarshal original response
	err := json.Unmarshal(searxJSONResponse, &resp)
	if err != nil {
		return nil, err
	}

	// Post processing of the results
	for i := range resp.Results {
		// Skip unrecognized content types
		if resp.Results[i].ImgFmt == "" {
			continue
		}

		// Base64 encode for easier transport e.g. when piping output to JQ
		imgSrcB64 := base64.StdEncoding.EncodeToString([]byte(resp.Results[i].ImgSrc))
		thumbSrcB64 := base64.StdEncoding.EncodeToString([]byte(resp.Results[i].ThumbnailSrc))
		resp.Results[i].ImgSrcB64 = imgSrcB64
		resp.Results[i].ThumbnailSrcB64 = thumbSrcB64

		results = append(results, resp.Results[i])
	}

	// Update response results
	resp.Results = results

	return &resp, nil
}
