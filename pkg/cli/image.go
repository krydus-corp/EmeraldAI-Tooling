// Package cli provides the Cobra CLI commands
/*
 * File: image.go
 * Project: cmd
 * File Created: Sunday, 5th April 2020 7:58:49 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Wednesday, 5th January 2022 12:06:07 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.com/krydus/emeraldai/emerald-tooling/pkg/image"
)

var (
	contentType bool
)

// imageCmd represents the imaging command
var imageCmd = &cobra.Command{
	Use:   "image <image path>",
	Short: "Formats an image",
	Long:  `Formats the content type or size of image.`,
	Run: func(cmd *cobra.Command, args []string) {

		streamInput, _ := cmd.Flags().GetBool("stream")
		replace, _ := cmd.Flags().GetBool("replace")
		workers, _ := cmd.Flags().GetInt("workers")
		height, _ := cmd.Flags().GetInt("height")
		width, _ := cmd.Flags().GetInt("width")
		format, _ := cmd.Flags().GetString("format")
		silent, _ := cmd.Flags().GetBool("silent")

		var contentType *string
		if format != "" {
			contentType = &format
		}

		// Check if any positional args supplied if not streaming input
		if !streamInput {
			if len(args) == 0 {
				log.Error("Non-streaming input with 0 length args; exiting")
				os.Exit(1)
			}
		}

		// Configure size conversion parameters
		var sizeParams *image.SizeConversionParams
		if height != 0 || width != 0 {
			sizeParams = &image.SizeConversionParams{Height: height, Width: width}
		}

		// Initialize the Imager
		imgr, err := image.NewImager(
			workers,
			contentType,
			sizeParams,
		)
		if err != nil {
			log.Errorf("Error initializing imager: %s", err.Error())
			os.Exit(1)
		}

		// Kick off Imager worker routines
		err = imgr.Start()
		if err != nil {
			log.Errorf("Error starting imager: %s", err.Error())
			os.Exit(1)
		}
		defer imgr.Stop()

		if streamInput {
			// Streaming input
			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					inpath := scanner.Text()

					imgBytes, err := ioutil.ReadFile(strings.TrimSpace(inpath))
					if err != nil {
						log.Infof("error processing input input=%s, error=%v", inpath, err)
						continue
					}

					imgr.InChan <- &image.Image{ImageBytes: imgBytes, OriginalFilepath: inpath}
				}
			}()
		} else {
			// Positional arg input
			inpath := args[0]
			imgBytes, err := ioutil.ReadFile(strings.TrimSpace(inpath))
			if err != nil {
				log.Infof("error processing input input=%s, error=%v", inpath, err)
				return
			}

			imgr.InChan <- &image.Image{ImageBytes: imgBytes, OriginalFilepath: inpath}
		}

		// Wait for all input to be processed
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		received := 0

		for {
			select {
			case f := <-imgr.OutChan:
				received++

				// Verbose output
				if received%10 == 0 {
					log.Infof("processed: %d", received)
				}

				// Check if unable to process
				if f.Err != nil {
					log.Errorf("error processing file: error=%s", f.Err.Error())
				} else {
					if err := ioutil.WriteFile(f.ProcessedFilepath, f.ImageBytes, 0666); err != nil {
						log.Errorf("error writing new image '%s'", path.Base(f.ProcessedFilepath))
					}

					if !silent {
						fmt.Fprint(os.Stdout, f.ProcessedFilepath+"\n")
					}

					if replace && f.ProcessedFilepath != f.OriginalFilepath {
						os.Remove(f.OriginalFilepath)
					}
				}

				// If not streaming, break
				if !streamInput {
					return
				}

			case <-sigs:
				log.Info("Received shutdown signal, exiting")
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)

	// Optional args
	imageCmd.Flags().IntP("workers", "w", 4, "Number of workers to process images")
	imageCmd.Flags().StringP("format", "f", "", "Format to standardize images")
	imageCmd.Flags().IntP("width", "x", 0, "Width to standardize image")
	imageCmd.Flags().IntP("height", "y", 0, "Height to standardize image")
	imageCmd.Flags().Bool("stream", false, "Streaming input")
	imageCmd.Flags().Bool("replace", false, "Replace original image")
	imageCmd.Flags().Bool("silent", false, "Do not output downloaded filepaths")
}
