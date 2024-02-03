/*
 * File: encode_test.go
 * Project: image
 * File Created: Sunday, 5th April 2020 3:36:47 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Wednesday, 5th January 2022 12:06:07 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package image

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	test_utils "gitlab.com/krydus/emeraldai/emerald-tooling/pkg/test"
)

func TestEncode(t *testing.T) {
	const (
		Width  = 400
		Height = 200
	)

	jpegImg, _ := test_utils.NewImage("image/jpeg", Width, Height)
	pngImg, _ := test_utils.NewImage("image/png", Width, Height)

	tests := map[string]struct {
		img        []byte
		conversion string
		err        bool
	}{
		"PNG to JPEG":         {pngImg, "image/jpeg", false},
		"JPEG to PNG":         {jpegImg, "image/png", false},
		"JPEG to Unsupported": {jpegImg, "image/svg", true},
	}

	for name, test := range tests {
		t.Logf("Running test %s", name)

		imgResult, err := ConvertImgType(test.img, test.conversion)
		if err != nil {
			// Not expecting an error on this test
			if !test.err {
				t.Errorf("failed converting image to MIME type '%s'; %v", test.conversion, err)
			}
			continue
		}

		mimeType := http.DetectContentType(imgResult)

		assert.Equal(t, test.conversion, mimeType)
	}
}
