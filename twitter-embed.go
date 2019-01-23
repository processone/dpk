package dpk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/processone/dpk/pkg/semweb"
)

func twitterEmbed(l Link) string {
	fmt.Println("Processing link:", l.URL)
	apiEndpoint := fmt.Sprintf("https://publish.twitter.com/oembed?url=%s", l.URL)
	client := semweb.NewClient()
	body, _, err := client.Get(apiEndpoint)
	if err != nil {
		fmt.Println(err)
		return l.Markdown()
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err)
		return l.Markdown()
	}

	var embed OEmbed
	if err = json.Unmarshal(data, &embed); err != nil {
		fmt.Println(err)
		return l.Markdown()
	}
	// Remove Javascript
	policy := bluemonday.UGCPolicy()
	policy.AllowStyling()
	safeHTML := policy.Sanitize(embed.HTML)

	// Parse HTML to rewrite links
	html2 := enrichHTML(safeHTML)

	return "\n" + html2
}

func enrichHTML(fragment string) string {
	// ======================================
	// Parsing
	context := html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	}
	doc, err := html.ParseFragment(strings.NewReader(fragment), &context)
	if err != nil {
		fmt.Println("html parsing error:", err)
		return fragment
	}

	// ======================================
	// Iterate on all elements in fragment (there is no necessarily a single Root)
	// to rewrite links
	for _, n := range doc {
		walkNode(n)
	}

	// ======================================
	// Render all elements from fragment.
	buf := new(bytes.Buffer)
	for _, n := range doc {
		if err := html.Render(buf, n); err != nil {
			// Return unmodified fragment on error
			fmt.Println("render error:", err)
			return fragment
		}
	}
	// Return modified fragment
	return buf.String()
}

func walkNode(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		var l Link
		for _, a := range n.Attr {
			if a.Key == "href" {
				l.URL = a.Val
				break
			}
		}

		// If we have a simple link content, we can rewrite it as well
		if c := n.FirstChild; c.Type == html.TextNode && c.NextSibling == nil {
			l.AnchorText = c.Data
			l = l.Resolve()
			if needRewrite(c.Data) {
				c.Data = l.URLTitle
			}

			// TODO: keep other attributes, like class
			var newAttr []html.Attribute
			n.Attr = append(newAttr, html.Attribute{"", "href", l.URL})
			return
		}

		// Otherwise, resolve only the URL and keep the children untouched
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walkNode(c)
		}

		l = l.Resolve()
		n.Attr = []html.Attribute{{Key: "href", Val: l.URL}}
		return
	}

	// Just walk down the node structure if this is not a node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNode(c)
	}
}

func needRewrite(urlText string) bool {
	if strings.HasPrefix(urlText, "http://") {
		return true
	}
	if strings.HasPrefix(urlText, "https://") {
		return true
	}
	if strings.HasPrefix(urlText, "pic.twitter.com/") {
		return true
	}
	return false
}
