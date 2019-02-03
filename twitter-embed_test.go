package dpk

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/processone/dpk/pkg/twitter"

	"github.com/processone/dpk/pkg/httpmock"
)

// TODO(mr): Test should be able to use mock: We need to be able to override the HTTP client.
func TestRewriteEmbedded(t *testing.T) {
	// Create a temp directory to store images
	dir, err := ioutil.TempDir("", "twitpic-test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // Clean up at end of test

	// Setup HTTP Mock
	client := httpmock.NewClient("fixtures/")
	fixtureName := "GetImage"
	if err := client.LoadScenario(fixtureName); err != nil {
		t.Errorf("Cannot load fixture %s: %s", fixtureName, err)
		return
	}

	expectedFilename := "DeywOwwWsAIij8t.jpg:large"
	example := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://t.co/wyyzANIGmB\" rel=\"nofollow\">https://t.co/wyyzANIGmB</a> <a href=\"https://t.co/OPYJZhQ9ih\" rel=\"nofollow\">pic.twitter.com/OPYJZhQ9ih</a></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\" rel=\"nofollow\">June 3, 2018</a></blockquote>"
	expected := fmt.Sprintf("<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://www.theverge.com/2018/6/3/17422752/microsoft-github-acquisition-rumors?utm_campaign=theverge&amp;utm_content=chorus&amp;utm_medium=social&amp;utm_source=twitter\">Microsoft has reportedly acquired GitHub</a> <img src=\"%s\"/></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%%5Etfw\">June 3, 2018</a></blockquote>",
		expectedFilename)
	result := enrichHTML(example, dir)

	// Check that we rewrite short links and image URLs
	if result != expected {
		t.Errorf("Result is not expect:\n%s", result)
	}

	// Check that the embedded image has been downloaded
	fullpath := filepath.Join(dir, expectedFilename)
	if _, err := os.Stat(fullpath); err != nil {
		t.Errorf("Image file was not downloaded: %s", fullpath)
	}
}

func TestGetImageUrl(t *testing.T) {
	// Setup HTTP Mock
	mock := httpmock.NewClient("fixtures/")
	// Scenario generated with:
	// httprec https://pic.twitter.com/ncJzTbz3dT GetImage
	// httprec https://pbs.twimg.com/media/DuIZsfQX4AAZFbs.png:large GetImage
	fixtureName := "GetImage"
	if err := mock.LoadScenario(fixtureName); err != nil {
		t.Errorf("Cannot load fixture %s: %s", fixtureName, err)
		return
	}

	picture := twitter.NewPicture()
	picture.SetBehaviour(mock)
	twitterPic := "https://pic.twitter.com/ncJzTbz3dT"
	imageUrl, err := picture.GetImageURL(twitterPic)
	if err != nil {
		t.Errorf("Cannot read page body: %s", err)
		return
	}

	if imageUrl != "https://pbs.twimg.com/media/DuIZsfQX4AAZFbs.png:large" {
		t.Errorf("Incorrect image target: %s", imageUrl)
	}
}
