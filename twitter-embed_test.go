package dpk

import (
	"testing"
)

// TODO: Mock HTTP requests
func TestRewriteEmbedded(t *testing.T) {
	example := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://t.co/wyyzANIGmB\" rel=\"nofollow\">https://t.co/wyyzANIGmB</a> <a href=\"https://t.co/OPYJZhQ9ih\" rel=\"nofollow\">pic.twitter.com/OPYJZhQ9ih</a></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\" rel=\"nofollow\">June 3, 2018</a></blockquote>"
	expected := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://www.theverge.com/2018/6/3/17422752/microsoft-github-acquisition-rumors?utm_campaign=theverge&amp;utm_content=chorus&amp;utm_medium=social&amp;utm_source=twitter\">Microsoft has reportedly acquired GitHub</a> <a href=\"https://twitter.com/verge/status/1003370586410815488/photo/1\">The Verge on Twitter</a></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\">June 3, 2018</a></blockquote>"
	result := enrichHTML(example)

	if result != expected {
		t.Errorf("Result is not expect:\n%s", result)
	}
}
