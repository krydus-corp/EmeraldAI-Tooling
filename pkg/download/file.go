// Package download provides file download utilities
/*
 * File: file.go
 * Project: download
 * File Created: Sunday, 22nd March 2020 7:25:52 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Tuesday, 7th April 2020 4:49:54 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package download

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/pkg/errors"
)

// File is a struct for representing the metadata of a saved file
type File struct {
	RawURL       string `json:"raw_url"`
	SanitizedURL string `json:"sanitized_url"`
	Name         string `json:"name"`
	TimeStamp    int64  `json:"timestamp"`
	Location     string `json:"location"`

	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
	FileBytes   []byte `json:"-"`

	Error error `json:"error"`
}

// newFile is a function for initializing a new file with pre-retrieval information
// Get will always return a non-nil result as any errors are encapsulated in the File object.
func newFile(urlStr, dst string) *File {
	f := &File{
		RawURL:    urlStr,
		TimeStamp: time.Now().Unix(),
		Location:  dst,
	}

	f.SanitizedURL, f.Error = validateURL(urlStr)
	if f.Error != nil {
		return f
	}

	u, err := stripQueryString(f.SanitizedURL)
	if err != nil {
		f.Error = err
		return f
	}

	f.Name = path.Base(u)

	return f
}

// get is a method for retrievingthe bytes of a file associated with a URL.
// get will always return a non-nil result as any errors are encapsulated
// in the File object.
//
// If File.Location is an empty string, the file will not be downloaded to the file system.
//
// get will block until an error or the file is received.
func (f *File) get() {
	resp, err := http.Get(f.SanitizedURL)
	if err != nil {
		f.Error = errors.Wrapf(err, "failed issuing a GET response for url [%s]", f.SanitizedURL)
		return
	}
	defer resp.Body.Close()

	fileBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errors.Wrapf(err, "failed reading GET response for url [%s]", f.SanitizedURL)
		return
	}

	if f.Location != "" {
		// Create the file
		name := path.Join(f.Location, f.Name)
		err := ioutil.WriteFile(name, fileBytes, 0666)
		if err != nil {
			f.Error = err
			return
		}
	}

	f.Size = len(fileBytes)
	f.FileBytes = fileBytes
	f.ContentType, err = getFileContentType(f.FileBytes)
	if err != nil {
		f.Error = err
		return
	}
}

// ToJSON is a function for exporting a File to JSON
// We do not include the actual file bytes if present when serializing to JSON.
func (f *File) ToJSON() ([]byte, error) {
	return json.MarshalIndent(f, "", "\t")
}

// getFileContentType is a helper function to retrieve the content type of a file
func getFileContentType(fileBytes []byte) (string, error) {
	reader := bytes.NewReader(fileBytes)

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := reader.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
