package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goversions",
	Short: "Get information about Go releases",
	Long:  `Get information about Go releases`,
}

func main() {
	rootCmd.AddCommand(listCommand())
	rootCmd.AddCommand(nextCommand())
	rootCmd.Execute()
}

var (
	allFlag bool
)

func listCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List Go releases",
		Run:     listAction,
	}

	cmd.Flags().BoolVarP(&allFlag, "all", "", false, "get all releases")

	return &cmd
}

func listAction(cmd *cobra.Command, args []string) {
	results := fetchReleases(allFlag)

	for _, result := range results {
		fmt.Println(result.Version)
	}
}

func nextCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "next",
		Aliases: []string{"l"},
		Short:   "Get next release (beta,rc)",
		Run:     nextAction,
	}

	return &cmd
}

func nextAction(cmd *cobra.Command, args []string) {
	results := fetchReleases(true)

	for _, result := range results {
		if strings.Index(result.Version, "beta") > -1 ||
			strings.Index(result.Version, "rc") > -1 {
			fmt.Println(result.Version)
		}
	}
}

func fetchReleases(all bool) []Result {
	url := "https://golang.org/dl/?mode=json"
	if all {
		url = "https://golang.org/dl/?mode=json&include=all"
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get %s: %s", url, err)
	}

	if resp.StatusCode != 200 {
		body := make([]byte, 1<<17)
		n, _ := io.ReadFull(resp.Body, body)

		log.Fatalf("Non 200 status=%d\n%s", resp.StatusCode, body[:n])
	}

	dec := json.NewDecoder(resp.Body)
	var results []Result
	err = dec.Decode(&results)
	if err != nil {
		log.Fatalf("Failed to parse releases: %s", err)
	}
	return results
}

type Result struct {
	Files []struct {
		Arch     string `json:"arch"`
		Filename string `json:"filename"`
		Kind     string `json:"kind"`
		Os       string `json:"os"`
		Sha256   string `json:"sha256"`
		Size     int64  `json:"size"`
		Version  string `json:"version"`
	} `json:"files"`
	Stable  bool   `json:"stable"`
	Version string `json:"version"`
}
