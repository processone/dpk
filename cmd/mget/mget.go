package main

import (
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

	title := getMetadata(args[0])
	fmt.Println(title)
}

func usage() {
	fmt.Println("Usage: mget [URL]")
}

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

func getMetadata(link string) string {
	title := ""
	client := httpClient()
Loop:
	for redirect := 0; redirect <= maxRedirect; redirect++ {
		resp, err := client.Get(link)
		if err != nil {
			fmt.Println(err)
			return title
		}

		switch resp.StatusCode {
		case 301, 302:
			location := resp.Header.Get("Location")
			fmt.Println("=> Resolved as", location)

			u, err := url.Parse(location)
			if err != nil {
				// Not a valid URL, just return the original link as is
				_ = resp.Body.Close()
				break Loop
			}
			// Retry resolving the next link, with new discovered value
			title = u.Host
			link = location
		case 200:
			page, err := metadata.ReadPage(resp.Body)
			if err != nil {
				fmt.Println("Cannot read metadata: ", err)
			} else {
				title = page.GetTitle()
			}
			_ = resp.Body.Close()
			break Loop
		default:
			_ = resp.Body.Close()
			break Loop
		}
	}

	return title
}
