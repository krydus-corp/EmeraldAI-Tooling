// Package cli provides the Cobra CLI commands
/*
 * File: download.go
 * Project: cmd
 * File Created: Tuesday, 24th March 2020 6:36:35 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Wednesday, 5th January 2022 12:06:07 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/krydus/emeraldai/emerald-tooling/pkg/download"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <URIs>",
	Short: "Download content from URI.",
	Long: `Download URIs using a raw, comma delimited string of URIs as input. 
	If the '--stream' option is supplied, the tool will run and accept URIs until terminated.
	Downloads content to the directory specified by --outpath (defaults to current directory).
	Outputs filepaths to STDOUT unless --silent option is specified.`,
	Run: func(cmd *cobra.Command, args []string) {

		outpath, _ := cmd.Flags().GetString("path")
		workers, _ := cmd.Flags().GetInt("workers")
		b64Encoded, _ := cmd.Flags().GetBool("b64")
		streamInput, _ := cmd.Flags().GetBool("stream")
		silent, _ := cmd.Flags().GetBool("silent")

		// Check if any positional args supplied if not streaming input
		urls := []string{}
		if !streamInput {
			if len(args) == 0 {
				log.Error("Non-streaming input with 0 length args; exiting")
				os.Exit(1)
			}
			urls = append(urls, strings.Split(args[0], ",")...)
		}

		// Initialize downloader
		downloader, err := download.NewDownloader(workers, outpath, false, b64Encoded)
		if err != nil {
			log.Errorf("Error initializing downloader: %s", err.Error())
			os.Exit(1)
		}

		// Kick off downloader worker routines
		err = downloader.Start()
		if err != nil {
			log.Errorf("Error starting downloader: %s", err.Error())
			os.Exit(1)
		}
		defer downloader.Stop()

		// Kick off an input streamer routine
		if streamInput {
			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					url := strings.TrimSpace(scanner.Text())
					downloader.InChan <- url
				}
			}()
		} else {
			// Queue up filepaths
			for _, url := range urls {
				downloader.InChan <- url
			}
		}

		// Wait for all sent filepaths to be processed
		sent := len(urls)
		received := 0

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case f := <-downloader.OutChan:
				received++

				// Verbose output
				if received%10 == 0 || sent == received {
					log.Infof("processed: %d/%d", received, sent)
				}

				// Check if unable to process
				if f.Error != nil {
					log.Errorf("error processing file: url=%s name=%s type=%s error=%s", f.RawURL, f.Name, f.ContentType, f.Error.Error())
				} else {
					log.Infof("downloaded file: url=%s name=%s type=%s size=%d\n", f.SanitizedURL, f.Name, f.ContentType, f.Size)
				}

				// Output filepath to stdout
				if !silent {
					fmt.Fprint(os.Stdout, path.Join(f.Location, f.Name)+"\n")
				}

				// Only triggers if went supplied URIs as a positional argument
				if sent == received {
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
	rootCmd.AddCommand(downloadCmd)

	// Optional args
	downloadCmd.Flags().StringP("path", "p", ".", "Path to output files")
	downloadCmd.Flags().IntP("workers", "w", 4, "Number of workers to process urls")
	downloadCmd.Flags().BoolP("b64", "b", false, "Base64 encoded input")
	downloadCmd.Flags().Bool("stream", false, "Streaming input")
	downloadCmd.Flags().Bool("silent", false, "Do not output downloaded filepaths")
}
