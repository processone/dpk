package main

import (
	"fmt"
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

// TODO(mr): Support method to be able to record posts.
// TODO(mr): Make it possible to record a sequence of several URL to generate more complex artefact.

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Missing url or name")
		os.Exit(1)
	}

	var seq httpmock.Sequence
	uri := args[0]
	seqName := args[1]
	client := newHTTPClient()

	// TODO(mr): Refactor to move to httpmock package.
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
			r.BodyFilename = fmt.Sprintf("%s-%d.html", seqName, redirect+1)
			step := httpmock.Step{
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

	filename := seqName + ".json"
	if err := seq.SaveTo(filename); err != nil {
		fmt.Printf("Cannot save sequence to file %s: %s", filename, err)
		os.Exit(2)
	}

	// Test ReadSequence:
	sequence, err := httpmock.ReadSequence(filename)
	if err != nil {
		fmt.Printf("Cannot read sequence from %s: %s", filename, err)
		os.Exit(3)
	}
	fmt.Println(sequence)
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
