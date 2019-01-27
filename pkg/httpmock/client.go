package httpmock

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

// HTTPMock is an HTTP client with a custom transport to simulate HTTP queries
type HTTPMock struct {
	*http.Client
	FixtureDir string

	count int
}

// Create a new HTTP client using a custom responder.
func NewClient(fixtureDir string) HTTPMock {
	// Create HTTPClient
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	c := http.Client{Timeout: time.Second * 15, Transport: transport}
	return HTTPMock{&c, fixtureDir, 0}
}

func (m *HTTPMock) LoadFixture(fixtureName string) error {
	m.count = 0
	filename := fixtureName + ".json"
	seq, err := ReadSequence(filepath.Join(m.FixtureDir, filename))
	if err != nil {
		return err
	}

	responder := func(req *http.Request) (*http.Response, error) {
		fmt.Printf("Request %d\n", m.count)

		if len(seq.Steps) <= m.count {
			return nil, errors.New(fmt.Sprintf("Unexpected step %d", m.count))
		}

		curStep := seq.Steps[m.count]
		if req.URL.String() != curStep.RequestURL {
			return nil, errors.New(fmt.Sprintf("step %d not matching requested URL %s. Expecting %s",
				m.count, req.URL.String(), curStep.RequestURL))
		}
		resp, err := curStep.Response.ToHTTPResponse()

		// Increment or reset step counter
		if m.count < len(seq.Steps) {
			m.count += 1
		} else {
			m.count = 0
		}
		return resp, err
	}
	m.setResponder(responder)
	return nil
}

func (m *HTTPMock) setResponder(responder Responder) {
	var t Transport
	t.Responder = responder
	m.Transport = &t
}

//=============================================================================
// Transport mock: Make it possible to define custom functions to generate
// local HTTP response for requesting client.

// Transport is an http.Transport mock.
type Transport struct {
	Responder Responder
}

// Responder is the signature of the responder function, which generate HTTP
// responses for client.
type Responder func(*http.Request) (*http.Response, error)

// Transport implements http.RoundTripper to act as an HTTP client mock.
func (m *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// If responder is missing, returns an error:
	if m.Responder == nil {
		resp := http.Response{}
		return &resp, errors.New("no responder found: you need to load a fixture")
	}

	// Otherwise, delegate response to that responder
	return m.Responder(req)
}
