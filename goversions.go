package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

	max := ParsedVersion{}
	for _, result := range results {
		if strings.Index(result.Version, "beta") > -1 ||
			strings.Index(result.Version, "rc") > -1 {
			parsed, err := parseVersion(result.Version)
			if err != nil {
				panic(err)
			}
			if max.Less(*parsed) {
				max = *parsed
			}
		}
	}
	fmt.Println(max.String())

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

var versionRegex = regexp.MustCompile(`\Ago(\d+)\.(\d+)(?:\.(\d+))?(?:rc(\d+)|beta(\d+))?\z`)

func parseVersion(s string) (*ParsedVersion, error) {
	groups := versionRegex.FindStringSubmatch(s)
	if len(groups) < 2 {
		return nil, fmt.Errorf("invalid go version")
	}

	var v ParsedVersion
	var err error
	v.Major, err = strconv.Atoi(groups[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version, this shouldn't be possible")
	}

	v.Minor, err = strconv.Atoi(groups[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version, this shouldn't be possible")
	}
	v.Patch, _ = strconv.Atoi(groups[3])
	v.RC, _ = strconv.Atoi(groups[4])
	v.Beta, _ = strconv.Atoi(groups[5])

	return &v, nil
}

type ParsedVersion struct {
	Major int
	Minor int
	Patch int
	RC    int
	Beta  int
}

// returns true if pv is less than other
func (pv ParsedVersion) Less(other ParsedVersion) bool {
	if pv.Major != other.Major {
		return pv.Major < other.Major
	}
	if pv.Minor != other.Minor {
		return pv.Minor < other.Minor
	}
	if pv.Patch != other.Patch {
		return pv.Patch < other.Patch
	}
	if pv.RC != other.RC {
		return pv.RC < other.RC
	}
	if pv.Beta != other.Beta {
		return pv.Beta < other.Beta
	}
	return false
}

func (pv ParsedVersion) String() string {
	out := fmt.Sprintf("go%d.%d", pv.Major, pv.Minor)
	if pv.Patch > 0 {
		out += fmt.Sprintf(".%d", pv.Patch)
	}
	if pv.RC > 0 {
		out += fmt.Sprintf("rc%d", pv.RC)
	}
	if pv.Beta > 0 {
		out += fmt.Sprintf("beta%d", pv.Beta)
	}
	return out
}
