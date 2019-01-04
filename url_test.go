package dpk_test

import (
	"testing"

	"github.com/processone/dpk"
)

func TestRedirectUrl(t *testing.T) {
	originalUrl := "https://donate.mozilla.org/"
	location := "/en-US/index.html"
	newUrl, err := dpk.RedirectUrl(originalUrl, location)
	if err != nil {
		t.Errorf("Could not properly generate full redirect URL: %s", err)
		return
	}
	expected := "https://donate.mozilla.org/en-US/index.html"
	if newUrl != expected {
		t.Errorf("Incorrect redirect URL. Got: '%s' Expected: '%s'", newUrl, expected)
	}
}
