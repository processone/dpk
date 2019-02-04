package httpmock_test

import (
	"io/ioutil"
	"testing"

	"github.com/processone/dpk/pkg/httpmock"
)

func TestMultipleRuns(t *testing.T) {
	// Setup HTTP Mock
	mock := httpmock.NewMock("fixtures/")
	// Scenario is generated with: httprec https://t.co/tprDWoN8vm TwoSteps
	fixtureName := "TwoSteps"
	if err := mock.LoadScenario(fixtureName); err != nil {
		t.Errorf("Cannot load fixture %s: %s", fixtureName, err)
		return
	}

	// We should be able to successfully execute the mock request several times without any error:
	for i := 0; i < 2; i++ {
		resp, err := mock.Get("https://t.co/tprDWoN8vm")
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

func TestMultipleSequences(t *testing.T) {
	// Setup HTTP Mock
	mock := httpmock.NewMock("fixtures/")
	// Scenario generated with:
	// httprec https://pic.twitter.com/ncJzTbz3dT scenario1
	// httprec https://pbs.twimg.com/media/DuIZsfQX4AAZFbs.png:large scenario1
	fixtureName := "scenario1"
	if err := mock.LoadScenario(fixtureName); err != nil {
		t.Errorf("Cannot load fixture %s: %s", fixtureName, err)
		return
	}
	// We can get HTML page
	resp, err := mock.Get("https://pic.twitter.com/ncJzTbz3dT")
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Incorrect status on HTML get: %d", resp.StatusCode)
	}

	// We can get the image
	resp, err = mock.Get("https://pbs.twimg.com/media/DuIZsfQX4AAZFbs.png:large")
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Incorrect status on HTML get: %d", resp.StatusCode)
	}
}
