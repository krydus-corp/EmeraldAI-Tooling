/*
 * File: resize_test.go
 * Project: image
 * File Created: Tuesday, 14th April 2020 7:55:47 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Wednesday, 5th January 2022 12:06:07 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
	test_utils "gitlab.com/krydus/emeraldai/emerald-tooling/pkg/test"
)

func TestResize(t *testing.T) {
	const (
		OriginalWidth  = 400
		OriginalHeight = 200
	)

	jpegImg, _ := test_utils.NewImage("image/jpeg", OriginalWidth, OriginalHeight)

	tests := map[string]struct {
		img        []byte
		convHeight int
		convWidth  int
		err        bool
	}{
		"JPEG Resize": {jpegImg, 100, 100, false},
	}

	for name, test := range tests {
		t.Logf("Running test %s", name)

		newImg, err := ResizeImage(test.img, test.convWidth, test.convHeight)
		if err != nil {
			t.Errorf("Unexpected error resizing image; error=%v", err)
			continue
		}

		stats, err := GetStats(newImg)
		if err != nil {
			t.Errorf("Unexpected error getting image stats; error=%v", err)
			continue
		}

		t.Logf("Resized image from (%d, %d) -> (%d, %d)", OriginalHeight, OriginalWidth, test.convHeight, test.convWidth)
		assert.Equal(t, stats.Height, test.convHeight)
		assert.Equal(t, stats.Width, test.convWidth)
	}
}
