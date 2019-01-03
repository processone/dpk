package metadata_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/processone/dpk/pkg/metadata"
)

func TestTitle(t *testing.T) {
	testFile := "fixtures/links-1.html"
	data, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("Cannot read file: %s", testFile)
	}

	page, err := metadata.ReadPage(bytes.NewReader(data))
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
	if page.Title() != expected {
		t.Errorf("Could not extract correct title from '%s'. Got: '%s' Expected: '%s'", testFile, page.Title(), expected)
	}
}

func ExamplePage_Title() {
	html := `<!DOCTYPE html>
<html lang="en">
<head prefix="og: http://ogp.me/ns#">
    <meta charset="utf-8"/>
    <meta property="og:title" content="Open Graph title" />
</head>
<body><p>This is a test page</p></body>
</html>`
	if page, err := metadata.ReadPage(strings.NewReader(html)); err != nil {
		fmt.Println(page.Title())
		// Output: Open Graph title
	}
}
