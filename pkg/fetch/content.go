// Package fetch provides fetching utilities for various content types
/*
 * File: content.go
 * Project: fetch
 * File Created: Saturday, 21st March 2020 4:10:26 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Friday, 1st May 2020 7:38:51 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package fetch

import "strings"

// AvailableContentTypes is a type for representing content types available for fetching
var AvailableContentTypes = map[string]bool{
	"images": true,
}

// CheckIsAvailableContentType is a function for checking if a content type is among the supported types
func CheckIsAvailableContentType(t string) bool {
	if _, ok := AvailableContentTypes[strings.ToLower(t)]; ok {
		return true
	}

	return false
}
