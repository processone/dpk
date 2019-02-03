package httpmock

import (
	"errors"
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

// LoadScenario reads a scenario from a file and set HTTPMock client responder.
func (m *HTTPMock) LoadScenario(scenarioName string) error {
	filename := scenarioName + ".json"
	scn, _, err := InitScenario(filepath.Join(m.FixtureDir, filename))
	if err != nil {
		return err
	}

	// TODO: Extract responder in a separate function
	responder := func(req *http.Request) (*http.Response, error) {
		key := key{method: req.Method, url: req.URL.String()}
		ref := scn.index[key]
		if !ref.exist {
			return nil, errors.New("request not found in scenario")
		}

		// TODO: Consistency checks to avoid out of range errors ?
		seq := scn.Sequences[ref.sequenceID]
		curStep := seq.Steps[ref.stepID]
		resp, err := curStep.Response.ToHTTPResponse()

		return resp, err
	}
	m.setResponder(responder)
	return nil
}

// setResponder is used to defined your own responder for an HTTPMock client.
// It can be used for cases for which you to not have a scenario file to load.
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
