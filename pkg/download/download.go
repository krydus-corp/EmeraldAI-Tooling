// Package download provides file download utilities
/*
 * File: download.go
 * Project: download
 * File Created: Sunday, 29th March 2020 5:00:08 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Thursday, 28th May 2020 8:54:34 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package download

import (
	"encoding/base64"
	"fmt"
	"os"
	"sync"

	guuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Downloader is a struct for holding a downloader's context variables
type Downloader struct {
	DestinationPath string
	NoStore         bool
	Base64Encoded   bool

	// Concurrency of this Downloader
	Concurrency int

	// Sync vars
	InChan   chan string
	OutChan  chan *File
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewDownloader is a function for initializing a new Downloader
// By default, the Downloader is configured to download files to the specified destination.
// If wanting to return the bytes of the files only, set the 'NoStore' property of the Downloader to 'True'
func NewDownloader(concurrency int, dst string, noStore, b64Encoded bool) (*Downloader, error) {
	if _, err := os.Stat(dst); os.IsNotExist(err) && !noStore {
		return nil, fmt.Errorf("destination path does not exist: '%s'", dst)
	}

	return &Downloader{
		Concurrency:     concurrency,
		DestinationPath: dst,
		NoStore:         noStore,
		Base64Encoded:   b64Encoded,
		InChan:          make(chan string, 100),
		OutChan:         make(chan *File, 100),
		stopChan:        make(chan struct{}, 1),
	}, nil
}

// Start kicks off the Downloader worker routines
func (d *Downloader) Start() error {
	// Start a fixed number of goroutines to read and digest urls.
	d.wg.Add(d.Concurrency)
	for i := 0; i < d.Concurrency; i++ {
		go func() {
			d.work()
			d.wg.Done()
		}()
	}

	return nil
}

// Stop the Downloader worker routines and wait for graceful exit
func (d *Downloader) Stop() {
	close(d.stopChan)
	d.wg.Wait()
}

func (d *Downloader) work() {
	id := guuid.New()
	log.Debugf("Starting Downloader worker routine [%s]\n", id.String())
	defer log.Debugf("Exiting Downloader worker routine [%s]\n", id.String())

	for {
		select {
		case urlStr := <-d.InChan:
			if d.Base64Encoded {
				u, err := base64.StdEncoding.DecodeString(urlStr)
				if err != nil {
					log.Warnf("unable to decode b64 encoded url [%s]", urlStr)
					continue
				}
				urlStr = string(u)
			}

			f := newFile(urlStr, d.DestinationPath)
			if f.Error == nil {
				f.get()
			}

			d.OutChan <- f

		case <-d.stopChan:
			return
		}
	}
}
