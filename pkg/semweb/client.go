package semweb

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

//=============================================================================
// Custom HTTP client
// Control timeouts, and redirect policy

// Client adds safer default values to Go HTTP client and provide control
// on redirect behaviour.
type Client struct {
	Client      *http.Client
	MaxRedirect int
	// TODO: Support debug logger
}

func NewClient() Client {
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
	return Client{Client: &client, MaxRedirect: 7}
}

// Get returns a web page reader, following a predefined number of redirects.
func (c Client) Get(url string) (io.ReadCloser, error) {
	for redirect := 0; redirect <= c.MaxRedirect; redirect++ {
		resp, err := c.Client.Get(url)
		if err != nil {
			return nil, err
		}

		switch {
		case resp.StatusCode >= 300 && resp.StatusCode < 400:
			// Redirect
			location := resp.Header.Get("Location")
			// Retry resolving the next link, with new discovered location
			url, err = formatRedirectUrl(url, location)
			_ = resp.Body.Close()
			if err != nil {
				// Not a valid URL, do not redirect further
				return nil, err
			}
			// TODO: Display using debug or verbose option
			// fmt.Println("=> Resolved as", location)
		case resp.StatusCode == 200:
			// Success
			return resp.Body, nil
		default:
			_ = resp.Body.Close()
			return nil, fmt.Errorf("unexpected response code %d", resp.StatusCode)
		}

	}
	return nil, fmt.Errorf("maximum number of redirects reached")
}

// Follow redirect and return final URL
func (c Client) FollowRedirect(currentUrl string) string {
Loop:
	// Try to resolve link N times, as sometimes you can find a chain of redirects before
	// reaching the canonical link.
	for redirect := 0; redirect <= c.MaxRedirect; redirect++ {
		resp, err := c.Client.Get(currentUrl)
		if err != nil {
			fmt.Println(err)
			return currentUrl
		}

		switch resp.StatusCode {
		case 301, 302:
			locationHdr := resp.Header.Get("Location")
			// Retry resolving the next link, with new discovered location
			link, err := formatRedirectUrl(currentUrl, locationHdr)
			if err != nil {
				// Not a valid URL, just return the original link as is
				_ = resp.Body.Close()
				break Loop
			}
			// TODO: Display using debug or verbose option
			// fmt.Println("=> Resolved as", link)

			_, err = url.Parse(link)
			if err != nil {
				// Not a valid URL, just return the original link as is
				_ = resp.Body.Close()
				break Loop
			}
			// Retry resolving the next link, with new discovered value
			currentUrl = locationHdr
		case 200:
			/* TODO: Refactor this in a method to get the page  body + final URL in a same call
			page, err := metadata.ReadPage(resp.Body)
			if err == nil {
				displayUrl = page.Title()
			}
			*/
			resp.Body.Close()
			break Loop
		default:
			fmt.Println("Ignored HTTP Status Code:", resp.StatusCode)
			resp.Body.Close()
			break Loop
		}
	}

	return currentUrl
}

// TODO: Should this method be on Context, taking only new link ?
func (c Client) ResolveReference(base, href string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

//=============================================================================
// HTTP request helpers

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
