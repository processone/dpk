package semweb_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/processone/dpk/pkg/httpmock"

	"github.com/processone/dpk/pkg/semweb"
)

// TODO: Rewrite, based on new mock package.
func TestFollowRedirect(t *testing.T) {
	targetSite := "https://process-one.net"
	responder := func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "t.co" {
			resp := RedirectResponse(targetSite)
			resp.Request = req
			return resp, nil
		}
		if req.URL.Host == "process-one.net" {
			resp := SimplePageResponse("Target Page Title")
			resp.Request = req
			return resp, nil
		}
		t.Errorf("unknown host: %s", req.Host)
		return nil, errors.New("unknown host")
	}

	mock := httpmock.NewMock("")
	mock.SetResponder(responder)
	c := semweb.NewClient()
	c.HTTPClient = mock.Client
	uri := c.FollowRedirect("https://t.co/shortURL")
	if uri != targetSite {
		t.Errorf("unexpected uri: %s", uri)
	}
}

// Simple basic HTTP responses

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
