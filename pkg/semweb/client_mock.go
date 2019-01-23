package semweb

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

func NewMockClient(responder Responder) Client {
	c := NewClient()

	var mock Transport
	mock.Responder = responder
	c.Client.Transport = &mock

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
