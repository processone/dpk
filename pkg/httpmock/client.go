package httpmock

import (
	"errors"
	"net"
	"net/http"
	"time"
)

// HTTPMock is an HTTP client with a custom transport to simulate HTTP queries
type HTTPMock struct {
	*http.Client
	count int
}

// Create a new HTTP client using a custom responder.
func NewMockClient(responder Responder) HTTPMock {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	c := http.Client{Timeout: time.Second * 15, Transport: transport}

	var mock Transport
	mock.Responder = responder
	c.Transport = &mock

	return HTTPMock{&c, 0}
}

func (m *HTTPMock) SetResponder(responder Responder) {
	var t Transport
	t.Responder = responder
	m.Transport = &t
}

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

// Responders
func ConnectionFailure(*http.Request) (*http.Response, error) {
	resp := http.Response{}
	return &resp, errors.New("no responder found")
}
