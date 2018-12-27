package main

import (
	"fmt"
	"github.com/processone/dpk"
	"os"
)

// This tool is used to convert data from your Twitter archive to a set of Markdown files.
// You can request your data from Twitter at this URL: https://twitter.com/settings/your_twitter_data
func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Missing argument.")
		usage()
		os.Exit(1)
	}

	if err := dpk.TwitterToMD(args[0], args[1]); err != nil {
		fmt.Println(err)
	}
}

func usage() {
	fmt.Println("Usage: twitter-to-md [TwitterArchiveDir] [OutputDir]")
}
