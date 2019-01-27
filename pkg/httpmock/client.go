package httpmock

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//=============================================================================
// Tools to mock HTTP client transport to simulate HTTP queries

// Responder is the signature of the clientMock creation function.
type Responder func(*http.Request) (*http.Response, error)

// Transport is an http.Transport mock.
type Transport struct {
	Responder Responder
}

// Transport implements http.RoundTripper
func (m *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Responder != nil {
		return m.Responder(req)
	}
	return ConnectionFailure(req)
}

func NewClient() *http.Client {
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
	return &client
}

// Create a new HTTP client using a custom responder
func NewMockClient(responder Responder) *http.Client {
	c := NewClient()

	var mock Transport
	mock.Responder = responder
	c.Transport = &mock

	return c
}

// Responders
func ConnectionFailure(*http.Request) (*http.Response, error) {
	resp := http.Response{}
	return &resp, errors.New("no responder found")
}

// Simple basic responses

func RedirectResponse(location string) *http.Response {
	status := 301
	reader := bytes.NewReader([]byte{})
	header := http.Header{}
	header.Add("Location", location)
	response := http.Response{
		Status:     strconv.Itoa(status),
		StatusCode: status,
		Body:       ioutil.NopCloser(reader),
		Header:     header,
	}
	return &response
}

func SimplePageResponse(title string) *http.Response {
	status := 200
	template := `<html>
<head><title>%s</title></head>
<body><h2>%s</h2></body>
</html>`
	page := fmt.Sprintf(template, title, title)
	reader := strings.NewReader(page)
	response := http.Response{
		Status:     strconv.Itoa(status),
		StatusCode: status,
		Body:       ioutil.NopCloser(reader),
		Header:     http.Header{},
	}
	return &response
}
