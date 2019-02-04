package semweb_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/processone/dpk/pkg/semweb"
)

// TODO: Rewrite, based on new mock package.
func TestFollowRedirect(t *testing.T) {
	targetSite := "https://process-one.net"
	responder := func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "t.co" {
			resp := semweb.RedirectResponse(targetSite)
			resp.Request = req
			return resp, nil
		}
		if req.URL.Host == "process-one.net" {
			resp := semweb.SimplePageResponse("Target Page Title")
			resp.Request = req
			return resp, nil
		}
		t.Errorf("unknown host: %s", req.Host)
		return nil, errors.New("unknown host")
	}

	c := semweb.NewMockClient(responder)
	uri := c.FollowRedirect("https://t.co/shortURL")
	if uri != targetSite {
		t.Errorf("unexpected uri: %s", uri)
	}
}
