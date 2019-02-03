package main

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/processone/dpk/pkg/httpmock"
)

// Fully record HTTP request interactions, including redirects
// The result can be used as an input for testing.
// It is used to help simulating test involving HTTP requests.
// The result is saved as a .gob file, a native format to serialize Go struct.
// It also save a .url file containing the original URL. It can be useful to regenerate
// the data if we need to change the structure of the Sequence.

// params:
// - dir: Fixture dir (default: fixtures)
// - name: Test fixture name
// - url: URL retrieve to record

// Example:
// go run pkg/httpmock/cmd/httprec/httprec.go https://pic.twitter.com/OPYJZhQ9ih test
// go run pkg/httpmock/cmd/httprec/httprec.go https://pbs.twimg.com/media/DeywOwwWsAIij8t.jpg:large test
// go run pkg/httpmock/cmd/httprec/httprec.go https://pic.twitter.com/ncJzTbz3dT scenario1
// go run pkg/httpmock/cmd/httprec/httprec.go https://pbs.twimg.com/media/DuIZsfQX4AAZFbs.png:large scenario1

// TODO(mr): Support method to be able to record http "POST" requests

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Missing url or name")
		os.Exit(1)
	}

	var seq httpmock.Sequence
	uri := args[0]
	scnName := args[1]
	client := newHTTPClient()

	scn, created, err := httpmock.InitScenario(scnName + ".json")
	if err != nil {
		fmt.Println("Cannot read scenario:", err)
		os.Exit(2)
	}
	if !created {
		fmt.Printf("Existing scenario: it contains %d sequences.\n", scn.Count())
	}

	// TODO(mr): Refactor to move to httpmock package core.
	// Record a new sequence
Loop:
	for redirect := 0; redirect <= 10; redirect++ {
		resp, err := client.Get(uri)

		// on error
		if err != nil {
			fmt.Println("Recording error:", err)
			step := httpmock.Step{
				RequestURL: uri,
				Response:   httpmock.NewResponse(resp),
				Err:        err.Error(),
			}
			seq.Steps = append(seq.Steps, step)
			break Loop
		}

		// Add step to sequence and record body file
		r := httpmock.NewResponse(resp)
		// On success
		switch {
		// Record and follow redirects steps:
		case resp.StatusCode >= 300 && resp.StatusCode < 400:
			// Redirect
			location := resp.Header.Get("Location")
			fmt.Println("Recording redirect to:", location)
			step := httpmock.Step{
				Method:     "GET",
				RequestURL: uri,
				Response:   r,
			}
			seq.Steps = append(seq.Steps, step)
			_ = resp.Body.Close()
			// Retry resolving the next link, with new discovered location
			uri, err = formatRedirectUrl(uri, location)
			if err != nil {
				// Not a valid URL, do not redirect further
				break Loop
			}
		default:
			fmt.Println("Recording response:", resp.StatusCode)
			// TODO(mr): Use response header to decide about the file extension
			r.BodyFilename = fmt.Sprintf("%s-%d-%d%s", scnName, scn.Count()+1, redirect+1, extension(r.Header))
			step := httpmock.Step{
				Method:     "GET",
				RequestURL: uri,
				Response:   r,
			}
			if err := step.SaveBody(resp.Body, r.BodyFilename); err != nil {
				fmt.Printf("Cannot save body file %s: %s\n", r.BodyFilename, err)
			}
			seq.Steps = append(seq.Steps, step)
			_ = resp.Body.Close()
			break Loop
		}
	}

	// Add Sequence to scenario
	if err = scn.AddSequence(seq); err != nil {
		fmt.Println("Cannot add sequence to scenario:", err)
		os.Exit(2)
	}

	// Save file
	filename := scnName + ".json"
	if err := scn.SaveTo(filename); err != nil {
		fmt.Printf("Cannot save sequence to file %s: %s", filename, err)
		os.Exit(3)
	}
}

//=============================================================================
// Helpers

// newHTTPClient returns an http.Client that does not automatically follow redirects.
func newHTTPClient() *http.Client {
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
	return &client
}

// formatRedirectUrl returns a valid full URL from an original URL and a "Location" Header.
// It support local redirection on same host.
func formatRedirectUrl(originalUrl, locationHeader string) (string, error) {
	newUrl, err := url.Parse(locationHeader)
	if err != nil {
		return "", err
	}

	// This is a relative URL, we need to use the host from original URL
	if newUrl.Host == "" && newUrl.Scheme == "" {
		oldUrl, err := url.Parse(originalUrl)
		if err != nil {
			return "", err
		}
		newUrl.Host = oldUrl.Host
		newUrl.Scheme = oldUrl.Scheme
	}
	return newUrl.String(), nil
}

func extension(header http.Header) string {
	content := header["Content-Type"]
	for _, typ := range content {
		if ext, err := mime.ExtensionsByType(typ); err == nil {
			// We return the longest known extension
			bestExtension := ""
			for _, e := range ext {
				if len(e) > len(bestExtension) {
					bestExtension = e
				}
			}
			if bestExtension != "" {
				return bestExtension
			}
		}
	}
	return ".data"
}
