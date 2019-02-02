package dpk

import (
	"testing"

	"github.com/processone/dpk/pkg/httpmock"
)

// TODO(mr): Mock HTTP requests. I still need to be able to add images and be able to record several requests sequence in a scenario.
func TestRewriteEmbedded(t *testing.T) {
	example := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://t.co/wyyzANIGmB\" rel=\"nofollow\">https://t.co/wyyzANIGmB</a> <a href=\"https://t.co/OPYJZhQ9ih\" rel=\"nofollow\">pic.twitter.com/OPYJZhQ9ih</a></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\" rel=\"nofollow\">June 3, 2018</a></blockquote>"
	expected := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://www.theverge.com/2018/6/3/17422752/microsoft-github-acquisition-rumors?utm_campaign=theverge&amp;utm_content=chorus&amp;utm_medium=social&amp;utm_source=twitter\">Microsoft has reportedly acquired GitHub</a> <img src=\"DeywOwwWsAIij8t.jpg:large\"/></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\">June 3, 2018</a></blockquote>"
	result := enrichHTML(example, "")

	if result != expected {
		t.Errorf("Result is not expect:\n%s", result)
	}
}

func TestGetImageUrl(t *testing.T) {
	// Setup HTTP Mock
	client := httpmock.NewClient("fixtures/")
	fixtureName := "GetImageUrl"
	if err := client.LoadFixture(fixtureName); err != nil {
		t.Errorf("Cannot load fixture %s: %s", fixtureName, err)
		return
	}

	twitterPic := "https://pic.twitter.com/OPYJZhQ9ih"
	imageUrl, err := GetImageURL(twitterPic)
	if err != nil {
		t.Errorf("Cannot read page body: %s", err)
		return
	}

	if imageUrl != "https://pbs.twimg.com/media/DeywOwwWsAIij8t.jpg:large" {
		t.Errorf("Incorrect image target: %s", imageUrl)
	}
}
