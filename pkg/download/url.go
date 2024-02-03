// Package download provides file download utilities
/*
 * File: url.go
 * Project: download
 * File Created: Sunday, 29th March 2020 9:10:19 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Tuesday, 7th April 2020 4:42:55 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package download

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

// validateURL is a helper function for validating a URL
// If no URL scheme is found, the url is returned with prefixed forward slashes
// removed. This is to prevent download errors related to commonly prefixed forward
// slashed returned from Searx e.g. //live.staticflickr.com/5338/6912241516_ba31a52ea0_c.jpg
func validateURL(urlStr string) (string, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return "", errors.Wrapf(err, "unabled to parse request URI [%s]", urlStr)
	}

	// If no scheme, assume http
	if u.Scheme == "" && strings.HasPrefix(urlStr, "//") {
		urlStr = "http:" + urlStr
	}

	// Now attempt to validate
	if validURL := govalidator.IsRequestURL(urlStr); !validURL {
		return "", fmt.Errorf("request URI is not a valid URL [%s]", urlStr)
	}

	return urlStr, nil
}

// stripQueryString is a helper function for stripping the query params from a url
// This is useful when trying to parse out a filename from a url
func stripQueryString(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	return u.String(), nil
}
