/*
 * File: file_test.go
 * Project: download
 * File Created: Saturday, 11th April 2020 7:36:37 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Monday, 19th July 2021 8:20:42 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package download

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	imgURL   = "https://via.placeholder.com/150"
	tmpStore = "/tmp/"
)

func TestFileGet(t *testing.T) {

	tests := map[string]struct {
		url string
		err bool
	}{
		"Image Get": {imgURL, false},
	}

	for name, test := range tests {
		t.Logf("Running test %s", name)

		f := newFile(test.url, tmpStore)
		f.get()

		assert.NoError(t, f.Error)
		assert.FileExists(t, path.Join(f.Location, f.Name))
	}
}
