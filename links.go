package dpk

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// TODO Make it more generic: Extract a struct from HTML containing various useful metadata info.
// Return struct + struct method to get title based on defined priorities og:title > twitter:title > title

// First, try to get the opengraph metadata (og:title) and then fallback to HTML title
// TODO also extract og:image. e.g.:
// <meta property="og:image" content="https://gigaom.com/wp-content/uploads/sites/1/2011/01/sonosgroup-804x516.jpg" />
// TODO extract dcterms.title
//  example: <meta name='dcterms.title' content='Amazon&#8217;s dead serious about the enterprise cloud' />
//  on: https://gigaom.com/2012/11/21/amazons-dead-serious-about-the-enterprise-cloud/
func GetTitle(body io.Reader, title string) string {
	tokenizer := html.NewTokenizer(body)

Loop:
	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err != io.EOF {
				// Print unexpected errors
				fmt.Println(err)
			}
			break Loop

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "meta":
				isTitle := false
				content := ""
				for _, attr := range token.Attr {
					// TODO Also check twitter:title as fallback
					if attr.Key == "property" && attr.Val == "og:title" {
						isTitle = true
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if isTitle == true {
					title = content
					break Loop
				}
			case "title":
				//the next token should be the page title
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					// Use page title but keep on searching of open graph title, which is often more accurate
					title = tokenizer.Token().Data
				}
			}
		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "head" {
				// We finished processing HTML head, no more metadata expected.
				break Loop
			}
		}
	}

	return title
}
