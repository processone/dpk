package metadata_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/processone/dpk/pkg/metadata"
)

func TestGetTitle(t *testing.T) {
	testFile := "fixtures/links-1.html"
	data, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("Cannot read file: %s", testFile)
	}

	page, err := metadata.FromReader(bytes.NewReader(data))
	if err != nil {
		t.Errorf("cannot read metadata: %v", err)
		return
	}

	// Properly extract Twitter title
	expected := "Twitter title"
	if page.Properties["twitter:title"] != expected {
		t.Errorf("Could not extract Twitter title from '%s'. Got: '%s' Expected: '%s'", testFile, page.Properties["twitter:title"], expected)
	}

	// Open Graph has a high priority
	expected = "Open Graph title"
	if page.GetTitle() != expected {
		t.Errorf("Could not extract correct title from '%s'. Got: '%s' Expected: '%s'", testFile, page.GetTitle(), expected)
	}
}
