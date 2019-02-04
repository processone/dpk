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
	HTTPClient  *http.Client
	MaxRedirect int
	// TODO: Support debug logger
}

type Response struct {
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
	}
	return Client{HTTPClient: &client, MaxRedirect: 7}
}

// Get returns a web page reader, following a predefined number of redirects.
// It also return the final URL
func (c Client) Get(url string) (io.ReadCloser, string, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, url, err
	}
	finalURL := resp.Request.URL.String()
	return resp.Body, finalURL, nil
}

// Follow redirect and return final URL
func (c Client) FollowRedirect(currentUrl string) string {
	resp, err := c.HTTPClient.Get(currentUrl)
	if err != nil {
		fmt.Println("MREMOND 1")
		return currentUrl
	}

	if resp.Request == nil {
		fmt.Println("MREMOND 2")
		return currentUrl
	}
	if resp.Request.URL == nil {
		fmt.Println("MREMOND 3")
		return currentUrl
	}

	finalURL := resp.Request.URL.String()
	return finalURL
}

// TODO: Should this method be on Context, taking only new link ?
func (c *Client) ResolveReference(base, href string) string {
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
