package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/processone/dpk/pkg/semweb"
)

// mget is a tools to extract standard information from page metadata, microformats, etc.
//
//
// mget can be used to extract page metadata, by just passing it a page URL.
//
// Usage:
//    mget https://www.process-one.net/
//
// mget can also be used for more advanced topic by using a specialized command.
//
// - `profiles`: mget can be use to get a view of all user identities starting from a profile page:
//
// Usage:
//    mget profiles [URL]

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Missing command or url")
		usage()
		os.Exit(1)
	}

	if len(args) == 1 {
		// Retrieve page and extract metadata
		getPageMetadata(args[0])
	}

	if len(args) >= 2 {
		command := args[0]
		switch command {
		case "profiles":
			err := getProfiles(args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("- Get page metadata as JSON")
	fmt.Println("  mget [URL]")
	fmt.Println("")
	fmt.Println("- Crawl pages from starting point to gather list of user profiles")
	fmt.Println("Usage: mget profiles [URL]")
}

//=============================================================================
// Page metadata command

func getPageMetadata(pageURL string) error {
	pageMeta, err := getMetadata(pageURL)
	if err != nil {
		fmt.Println("Cannot retrieve page metadata:", err)
		os.Exit(1)
	}
	// Convert page info to JSON
	jsonMeta, err := json.MarshalIndent(pageMeta, "", "\t")
	if err != nil {
		fmt.Println("Cannot serialize metadata:", err)
		os.Exit(1)
	}
	// Prints JSON on stdout
	fmt.Println(string(jsonMeta))
	return nil
}

func getMetadata(link string) (semweb.Page, error) {
	var page semweb.Page
	client := semweb.NewClient()
	body, _, err := client.Get(link)
	if err != nil {
		return page, err
	}
	defer body.Close()

	page, err = semweb.ReadPage(body)
	if err != nil {
		return page, err
	}
	return page, nil
}

//=============================================================================
// Profile crawler

func getProfiles(profileURL string) error {
	client := semweb.NewClient()
	body, _, err := client.Get(profileURL)
	if err != nil {
		return err
	}
	defer body.Close()

	// TODO: extract profile info and add unknown URLs to the list of discovered profiles
	// Be careful:  We only need to keep bidirectionally certified profiles to avoid spammy URL
	// Probably we can return a list of certified profile, separated by a list of possible risky profile (we will not
	// crawl them further).
	ctx := semweb.Context{Client: client, Url: profileURL}
	urls, err := semweb.ExtractRelMe(ctx, body)
	if err != nil {
		return err
	}

	for _, u := range urls {
		if targetUrl := client.FollowRedirect(u); targetUrl != "" {
			fmt.Println(targetUrl)
		}
	}

	return nil
}

/*
TODO:

- Add -o option to output to a file
- Support Turtle syntax of metadata output ?
- Support Text output ?

*/
