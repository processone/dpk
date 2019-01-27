package dpk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/processone/dpk/pkg/httpmock"
)

// TODO: Mock HTTP requests
func TestRewriteEmbedded(t *testing.T) {
	example := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://t.co/wyyzANIGmB\" rel=\"nofollow\">https://t.co/wyyzANIGmB</a> <a href=\"https://t.co/OPYJZhQ9ih\" rel=\"nofollow\">pic.twitter.com/OPYJZhQ9ih</a></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\" rel=\"nofollow\">June 3, 2018</a></blockquote>"
	expected := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://www.theverge.com/2018/6/3/17422752/microsoft-github-acquisition-rumors?utm_campaign=theverge&amp;utm_content=chorus&amp;utm_medium=social&amp;utm_source=twitter\">Microsoft has reportedly acquired GitHub</a> <img src=\"https://pbs.twimg.com/media/DeywOwwWsAIij8t.jpg:large\"/></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\">June 3, 2018</a></blockquote>"
	result := enrichHTML(example, "")

	if result != expected {
		t.Errorf("Result is not expect:\n%s", result)
	}
}

func TestGetImageUrl(t *testing.T) {
	fixtureName := "GetImageUrl"
	// TODO(mr): Extract mock and Sequence Responder as a standard httpmock helper:
	requestNumber := 0
	responder := func(req *http.Request) (*http.Response, error) {
		fmt.Printf("Request %d\n", requestNumber)
		seq, err := httpmock.ReadSequence("fixtures/" + fixtureName + ".json")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot read fixture %s: %s", fixtureName, err))
		}

		if len(seq.Steps) <= requestNumber {
			return nil, errors.New(fmt.Sprintf("Unexpected step %d", requestNumber))
		}

		curStep := seq.Steps[requestNumber]
		if req.URL.String() != curStep.RequestURL {
			return nil, errors.New(fmt.Sprintf("step %d not matching requested URL %s. Expecting %s",
				requestNumber, req.URL.String(), curStep.RequestURL))
		}
		resp, err := curStep.Response.ToHTTPResponse()
		requestNumber += 1
		return resp, err
	}

	//
	client := httpmock.NewMockClient(responder)
	resp, err := client.Get("https://pic.twitter.com/OPYJZhQ9ih")
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Cannot read page body: %s", err)
		return
	}
	//fmt.Printf("%s", page)
}
