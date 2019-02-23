package httpmock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Step is used to record one of the HTTP steps in a specific query.
type Step struct {
	RequestURL string
	Method     string
	Response   Response
	Err        string
}

// SaveBody writes the content of the body to a separate file during the recording of a query.
func (step Step) SaveBody(content io.Reader, toFile string) (err error) {
	file, err := os.Create(toFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

// Response is a simplified http.Response version tailored for serialization.
type Response struct {
	Status       string
	StatusCode   int
	Proto        string
	ProtoMajor   int
	ProtoMinor   int
	Header       http.Header
	BodyFilename string
}

// NewResponse generates a simple response from an http.Response.
func NewResponse(resp *http.Response) Response {
	var r Response

	if resp != nil {
		r.Status = resp.Status
		r.StatusCode = resp.StatusCode
		r.Proto = resp.Proto
		r.ProtoMajor = resp.ProtoMajor
		r.ProtoMinor = resp.ProtoMinor
		r.Header = resp.Header
	}
	return r
}

// ToHTTPResponse generates an http.Response with the data from the simple response.
func (r Response) ToHTTPResponse() (*http.Response, error) {
	var resp http.Response
	resp.Status = r.Status
	resp.StatusCode = r.StatusCode
	resp.Proto = r.Proto
	resp.ProtoMajor = r.ProtoMajor
	resp.ProtoMinor = r.ProtoMinor
	resp.Header = r.Header

	if r.BodyFilename != "" {
		// TODO(mr): Fixture dir should not be hardcoded
		file, err := os.Open("fixtures/" + r.BodyFilename) // TODO(mr): Use filename join for OS independance
		if err != nil {
			return &resp, err
		}
		resp.Body = file
	} else {
		reader := bytes.NewReader([]byte{})
		resp.Body = ioutil.NopCloser(reader)
	}
	return &resp, nil
}

// Sequence is used to record all the redirect steps of a single HTTP Response.
type Sequence struct {
	Steps []Step
}

type key struct {
	method string
	url    string
}

type ref struct {
	exist      bool
	sequenceID int
	stepID     int
}

type index map[key]ref

type Scenario struct {
	Sequences []Sequence
	// internal index to access a sequence to send the right replies
	index index
	path  string
}

func (scn *Scenario) Count() int {
	return len(scn.Sequences)
}

func (scn *Scenario) AddSequence(seq Sequence) error {
	// If one of the sequence URL is already in scenario, reject the sequence
	for _, step := range seq.Steps {
		key := key{method: step.Method, url: step.RequestURL}
		if ref := scn.index[key]; ref.exist == true {
			return fmt.Errorf("duplicate method %s for URL %s", step.Method, step.RequestURL)
		}
	}

	scn.Sequences = append(scn.Sequences, seq)
	scn.updateIndex()
	return nil
}

func (scn *Scenario) updateIndex() {
	index := make(map[key]ref)
	for seqI, seq := range scn.Sequences {
		for stepI, step := range seq.Steps {
			key := key{
				method: step.Method,
				url:    step.RequestURL,
			}
			index[key] = ref{
				exist:      true,
				sequenceID: seqI,
				stepID:     stepI,
			}
		}
	}
	scn.index = index
}

// SaveTo stores JSON data for a scenario.
func (scn *Scenario) SaveTo(filePath string) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(scn)
}

// SaveAsURLList stores a pure text file containing all the initial URLs used to generate a scenario
// This can be used to rebuild the scenario file if the file format changes.
func (scn *Scenario) SaveAsURLList(filePath string) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, seq := range scn.Sequences {
		if len(seq.Steps) > 0 {
			if _, err := file.WriteString(seq.Steps[0].RequestURL + "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

// InitScenario reads a JSON representation from a given file if it does exists
// and generate a valid scenario.
// If the file does not exist, it will create and empty scenario
// TODO: If file does not exists we create it, but we need to also takes into account and return other error cases.
func InitScenario(filePath string) (scn *Scenario, created bool, err error) {
	var s Scenario
	s.path = filePath
	file, err := os.Open(filePath)
	if err != nil {
		return &s, true, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	if err != nil {
		return &s, false, err
	}
	s.updateIndex()
	return &s, false, nil
}

// TODO: Rewrite 'date' and 'expires' headers when replaying the scenario
