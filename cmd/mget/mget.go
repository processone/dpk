package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/processone/dpk/pkg/metadata"
)

// Extract page title
// This is the equivalent of wget, but focus on extracting / getting page info
//
// Usage:
//    mget https://www.process-one.net/

var maxRedirect = 7

// TODO: Extract HTML page metadata as json
func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Missing url")
		usage()
		os.Exit(1)
	}

	// Retrieve page and extract metadata
	pageMeta, err := getMetadata(args[0])
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
}

func usage() {
	fmt.Println("Usage: mget [URL]")
}

func getMetadata(link string) (metadata.Page, error) {
	var page metadata.Page
	client := httpClient()
Loop:
	for redirect := 0; redirect <= maxRedirect; redirect++ {
		resp, err := client.Get(link)
		if err != nil {
			return page, err
		}

		switch resp.StatusCode {
		case 301, 302:
			location := resp.Header.Get("Location")
			fmt.Println("=> Resolved as", location)

			_, err := url.Parse(location)
			if err != nil {
				// Not a valid URL, just return the original link as is
				_ = resp.Body.Close()
				break Loop
			}
			// Retry resolving the next link, with new discovered location
			link = location
		case 200:
			page, err = metadata.ReadPage(resp.Body)
			if err != nil {
				return page, err
			}
			_ = resp.Body.Close()
			break Loop
		default:
			_ = resp.Body.Close()
			break Loop
		}
	}

	return page, nil
}

//=============================================================================
// HTTP client

// httpClient adds safer default values to Go HTTP client.
func httpClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Second * 15,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

/*
TODO:

- Add -o option to output to a file
- Support Turtle syntax ?

*/
