// Package image provides image processing utilities
/*
 * File: resize.go
 * Project: imaging
 * File Created: Saturday, 4th April 2020 10:48:09 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Sunday, 9th May 2021 12:20:23 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package image

import (
	"bytes"
	"fmt"
	"image"
	"net/http"

	"github.com/disintegration/imaging"
)

// ResizeImage is a function for resizing an image.
// If one of width or height is 0, the image aspect ratio is preserved.
func ResizeImage(imgBytes []byte, width, height int) ([]byte, error) {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}

	// Resize the image
	dstImg := imaging.Resize(img, width, height, imaging.Lanczos)

	// Encode back to original format
	var encodedImg *bytes.Buffer

	contentType := http.DetectContentType(imgBytes)
	switch contentType {
	case ContentTypeJPEG:
		encodedImg, err = encodeImageToJPEG(dstImg)
	case ContentTypePNG:
		encodedImg, err = encodeImageToPNG(dstImg)
	default:
		return nil, fmt.Errorf("unsupported MIME type '%s'", contentType)
	}

	if err != nil {
		return nil, err
	}

	return encodedImg.Bytes(), nil
}
