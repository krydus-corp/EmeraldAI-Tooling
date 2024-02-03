/*
 * File: thumbnail.go
 * Project: image
 * File Created: Sunday, 9th May 2021 12:22:42 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Sunday, 9th May 2021 12:38:55 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"net/http"

	"github.com/disintegration/imaging"
)

// Thumbnail is a function that scales the image up or down using the specified resample filter,
// crops it to the specified width and hight and returns a base64 encoded string of the transformed image.
func Thumbnail(imgBytes []byte, width, height int) (string, int, error) {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return "", 0, err
	}

	// Resize the image
	dstImg := imaging.Thumbnail(img, width, height, imaging.Lanczos)

	// Encode back to original format
	var encodedImg *bytes.Buffer

	contentType := http.DetectContentType(imgBytes)
	switch contentType {
	case ContentTypeJPEG:
		encodedImg, err = encodeImageToJPEG(dstImg)
	case ContentTypePNG:
		encodedImg, err = encodeImageToPNG(dstImg)
	default:
		return "", 0, fmt.Errorf("unsupported MIME type '%s'", contentType)
	}

	if err != nil {
		return "", 0, err
	}

	sEnc := base64.StdEncoding.EncodeToString(encodedImg.Bytes())

	return sEnc, len([]byte(sEnc)), nil
}
