// Package cli provides the Cobra CLI commands
/*
 * File: fetch.go
 * Project: cmd
 * File Created: Sunday, 22nd March 2020 1:40:10 pm
 * Author: krydus (krydus@proton.me)
 * -----
 * Last Modified: Wednesday, 5th January 2022 12:06:07 pm
 * Modified By: krydus (krydus@proton.me>)
 */
package cli

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.com/krydus/emeraldai/emerald-tooling/pkg/fetch"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch <query>",
	Short: "Fetch content from the metasearch engine.",
	Long: `Fetches content from a Searx instance using an internal HTTP connector. 
	
By default, the 'fetch' command attempts to connect to a Searx instance running locally on 
127.0.0.1:8080. If a Searx instance is not running locally, you can modify the connection endpoint
to point to a publically hosted instance. 
E.g. './emerald-cli fetch cat -t images -l english -u https://searx.info/'
The output of this command is the returned JSON of the content results to stdout. 
Simply pipe the results to a file is desired:  './emerald-cli fetch cat -t images > out.json

Pro Tip: To quickly check the number of results, use JQ: cat out.json | jq '.resultno'
Note that the 'resultno' field is only available when batch downloading i.e. not using the 'stream' option`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		query := args[0]

		server, _ := cmd.Flags().GetString("server")
		contentType, _ := cmd.Flags().GetString("type")
		stream, _ := cmd.Flags().GetBool("stream")
		pageno, _ := cmd.Flags().GetInt("pageno")
		pages, _ := cmd.Flags().GetInt("pages")
		languagesAllExt, _ := cmd.Flags().GetBool("all-langs-ext")
		languagesAllSimple, _ := cmd.Flags().GetBool("all-langs-simple")
		languages, err := cmd.Flags().GetStringSlice("languages")

		f, err := fetch.NewFetcher(server, contentType, languagesAllExt, languagesAllSimple, stream, languages...)
		if err != nil {
			log.Errorf("Error initializing Fetcher; %s", err.Error())
			os.Exit(1)
		}

		for i := 0; i < pages; i++ {
			res := f.FetchAsync(query, pageno)
			for {
				if res.Ready {
					break
				}

				time.Sleep(200 * time.Millisecond)
			}

			if res.HasErrors() {
				log.Warn("Fetch results contain errors; some or all results may not be present")
				log.Warn(res.Errors)
			}

			// Output result set if not streaming
			if !stream {
				jsonBytes, err := res.ToJSON()
				if err != nil {
					log.Errorf("unable to output fetch results to JSON; %s", err.Error())
					os.Exit(1)
				}

				// Output to stdout
				fmt.Fprint(os.Stdout, string(jsonBytes))
			}

			pageno++
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Required args
	fetchCmd.Flags().StringP("type", "t", "", "Content type (required)")
	fetchCmd.MarkFlagRequired("type")

	// Optional args
	fetchCmd.Flags().StringP("server", "s", "http://127.0.0.1:8080", "Searx server URL")
	fetchCmd.Flags().Bool("stream", false, "Stream results as they come in")
	fetchCmd.Flags().Bool("all-langs-ext", false, "Search all extended languages")
	fetchCmd.Flags().Bool("all-langs-simple", true, "Search all simplified languages")
	fetchCmd.Flags().StringSliceP("languages", "l", []string{}, "Languages to search in")
	fetchCmd.Flags().IntP("pageno", "p", 1, "Page number to search")
	fetchCmd.Flags().IntP("pages", "n", 1, "Pages to fetch")

}
