package httpmock_test

import (
	"io/ioutil"
	"testing"

	"github.com/processone/dpk/pkg/httpmock"
)

func TestSequenceCounterReset(t *testing.T) {
	// Setup HTTP Mock
	client := httpmock.NewClient("fixtures/")
	fixtureName := "TwoSteps"
	if err := client.LoadFixture(fixtureName); err != nil {
		t.Errorf("Cannot load fixture %s: %s", fixtureName, err)
		return
	}

	// We should be able to successfully execute the mock request twice.
	for i := 0; i < 2; i++ {
		resp, err := client.Get("https://t.co/tprDWoN8vm")
		if err != nil {
			t.Errorf("Get error: %s", err)
			return
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Cannot read page body: %s", err)
			return
		}
	}
}
