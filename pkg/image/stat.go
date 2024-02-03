/*
 * File: stat.go
 * Project: image
 * File Created: Saturday, 11th April 2020 7:25:31 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Saturday, 30th May 2020 5:41:19 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package image

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"io"
)

// Stats is a struct for representing an image's stats
type Stats struct {
	Height      int
	Width       int
	ContentType string
}

// GetStats is a function for retrieving an image's stats e.g. height, width
func GetStats(imgBytes []byte) (*Stats, error) {
	r := bytes.NewReader(imgBytes)
	image, fmt, err := image.DecodeConfig(r) // Image Struct
	if err != nil {
		return nil, err
	}

	return &Stats{
		Height:      image.Height,
		Width:       image.Width,
		ContentType: fmt,
	}, nil
}

// Hash is a function for hashing a byte slice and returning the MD5 checksum
func Hash(imageBytes []byte) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
