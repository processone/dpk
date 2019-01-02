package dpk

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestGetTitle(t *testing.T) {
	testFile := "fixtures/links-1.html"
	data, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("Cannot read file: %s", testFile)
	}

	result := GetTitle(bytes.NewReader(data), "Default title")
	expected := "Open Graph title"
	if result != expected {
		t.Errorf("Could not extract correct title from '%s'. Got: '%s' Expected: '%s'", testFile, result, expected)
	}
}
