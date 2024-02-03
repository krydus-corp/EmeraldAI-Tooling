// Package fetch provides fetching utilities for various content types
/*
 * File: lang.go
 * Project: fetch
 * File Created: Saturday, 21st March 2020 3:08:54 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Sunday, 29th March 2020 5:28:44 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package fetch

import (
	"strings"
)

// Languages is a map for representing languages available for fetching
var Languages = map[string]string{
	"arabic":        "ar-sa",
	"bulgarian":     "bg-bg",
	"catalan":       "ca-es",
	"czech":         "cs-cz",
	"danish":        "da-dk",
	"german":        "de",
	"german-at":     "de-at",
	"german-ch":     "de-ch",
	"german-de":     "de-de",
	"greek":         "el-gr",
	"english":       "en",
	"english-au":    "en-au",
	"english-ca":    "en-ca",
	"english-gb":    "en-gb",
	"english-in":    "en-in",
	"english-my":    "en-my",
	"english-us":    "en-us",
	"spanish":       "es",
	"spanish-ar":    "es-ar",
	"spanish-es":    "es-es",
	"spanish-mx":    "es-mx",
	"estonian":      "et-ee",
	"persian":       "fa-ir",
	"finnish":       "fi-fi",
	"french":        "fr",
	"french-be":     "fr-be",
	"french-ca":     "fr-ca",
	"french-ch":     "fr-ch",
	"french-fr":     "fr-fr",
	"hebrew":        "he-il",
	"croatian":      "hr-hr",
	"hungarian":     "hu-hu",
	"indonesian":    "id-id",
	"icelandic":     "is-is",
	"italian":       "it-it",
	"japanese":      "ja-jp",
	"korean":        "ko-kr",
	"lithuanian":    "lt-lt",
	"latvian":       "lv-lv",
	"malay":         "ms-my",
	"norwegian":     "nb-no",
	"dutch":         "nl",
	"dutch-be":      "nl-be",
	"dutch-nl":      "nl-nl",
	"polish":        "pl-pl",
	"portuguese":    "pt",
	"portuguese-br": "pt-br",
	"portuguese-pt": "pt-pt",
	"romanian":      "ro-ro",
	"russian":       "ru-ru",
	"slovak":        "sk-sk",
	"slovenian":     "sl-si",
	"serbian":       "sr-rs",
	"swedish":       "sv-se",
	"thai":          "th-th",
	"turkish":       "tr-tr",
	"ukrainian":     "uk-ua",
	"vietnamese":    "vi-vn",
	"chinese":       "zh",
	"chinese-cn":    "zh-cn",
	"chinese-tw":    "zh-tw",
}

// getLanguageCode returns the language codes associated with the passed in language string
func getLanguageCode(lang string) (string, bool) {
	lang = strings.ToLower(lang)

	if code, ok := Languages[lang]; ok {
		return code, true
	}

	return "", false
}

// getSimplifiedLanguageCodes returns the simplified language codes from all aviailable language codes
func getSimplifiedLanguageCodes() []string {
	var m []string

	for _, v := range Languages {
		if !strings.ContainsAny(v, "-") {
			m = append(m, v)
		}
	}

	return m
}

// getAllLanguageCodes returns the all available language codes
func getAllLanguageCodes() []string {
	var m []string

	for _, v := range Languages {
		m = append(m, v)
	}

	return m
}
