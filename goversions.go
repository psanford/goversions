package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goversions",
	Short: "Get information about Go releases",
	Long:  `Get information about Go releases`,
}

func main() {
	rootCmd.AddCommand(listCommand())
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
	url := "https://golang.org/dl/?mode=json"
	if allFlag {
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

	seen := make(map[string]bool)
	for _, result := range results {
		for _, f := range result.Files {
			if seen[f.Version] {
				continue
			}
			fmt.Println(f.Version)
			seen[f.Version] = true
		}
	}
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
