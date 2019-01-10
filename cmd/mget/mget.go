package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/processone/dpk"
	"github.com/processone/dpk/pkg/metadata"
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

var maxRedirect = 7

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

func getMetadata(link string) (metadata.Page, error) {
	var page metadata.Page
	client := newHttpClient()
	body, err := client.get(link, maxRedirect)
	if err != nil {
		return page, err
	}
	defer body.Close()

	page, err = metadata.ReadPage(body)
	if err != nil {
		return page, err
	}
	return page, nil
}

//=============================================================================
// profile crawler

type Profile struct {
	Description string
	URL         string
}
type Profiles []Profile

func getProfiles(profileURL string) error {
	client := newHttpClient()
	body, err := client.get(profileURL, maxRedirect)
	if err != nil {
		return err
	}
	defer body.Close()

	// TODO: extract profile info and add unknown URLs to the list of discovered profiles
	// Be careful:  We only need to keep bidirectionally certified profiles to avoid spammy URL
	// Probably we can return a list of certified profile, separated by a list of possible risky profile (we will not
	// crawl them further).

	return nil
}

//=============================================================================
// Custom HTTP client
// Control timeouts, and redirect policy

type httpClient struct {
	client *http.Client
	// TODO: Support debug logger
}

// httpClient adds safer default values to Go HTTP client.
func newHttpClient() httpClient {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	client := http.Client{
		Timeout:   time.Second * 15,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return httpClient{client: &client}
}

// get returns a web page, following a predefined number of redirects.
func (c httpClient) get(url string, maxRedirect int) (io.ReadCloser, error) {
	for redirect := 0; redirect <= maxRedirect; redirect++ {
		resp, err := c.client.Get(url)
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode {
		case 301, 302:
			location := resp.Header.Get("Location")
			// Retry resolving the next link, with new discovered location
			url, err = dpk.RedirectUrl(url, location)
			resp.Body.Close()
			if err != nil {
				// Not a valid URL, do not redirect further
				return nil, err
			}
			// TODO: Display using debug or verbose option
			// fmt.Println("=> Resolved as", location)
		case 200:
			return resp.Body, nil
		default:
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected response code %d", resp.StatusCode)
		}

	}
	return nil, fmt.Errorf("maximum number of redirect reached")
}

/*
TODO:

- Add -o option to output to a file
- Support Turtle syntax of metadata output ?
- Support Text output ?

*/
