package dpk

import (
	"fmt"
	"testing"
)

func TestRewriteEmbedded(t *testing.T) {
	example := "<blockquote class=\"twitter-tweet\"><p lang=\"en\" dir=\"ltr\">Microsoft has reportedly acquired GitHub <a href=\"https://t.co/wyyzANIGmB\" rel=\"nofollow\">https://t.co/wyyzANIGmB</a> <a href=\"https://t.co/OPYJZhQ9ih\" rel=\"nofollow\">pic.twitter.com/OPYJZhQ9ih</a></p>— The Verge (@verge) <a href=\"https://twitter.com/verge/status/1003370586410815488?ref_src=twsrc%5Etfw\" rel=\"nofollow\">June 3, 2018</a></blockquote>"
	result := enrichHTML(example)
	fmt.Println(result)
}
