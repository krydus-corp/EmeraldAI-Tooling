/*
 * File: main_test.go
 * Project: imaging_test
 * File Created: Sunday, 5th April 2020 3:44:51 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Saturday, 11th April 2020 11:30:11 am
 * Modified By: krydus (krydus@proton.me>)
 */
package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

// NewImage is a function for providing a test image
func NewImage(mimeType string, width, height int) ([]byte, error) {
	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw a red dot at (2, 3)
	img.Set(2, 3, color.RGBA{255, 0, 0, 255})

	var b = new(bytes.Buffer)
	var err error

	switch mimeType {
	case "image/jpeg":
		err = jpeg.Encode(b, img, nil)
	case "image/png":
		err = png.Encode(b, img)
	default:
		return nil, fmt.Errorf("unsupported MIME type '%s'", mimeType)
	}

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
