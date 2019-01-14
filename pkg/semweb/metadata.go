// It is more general and handle basic metadata, microformats, RDFa, etc.
package semweb

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Properties is a map gathering HTML page metadata properties.
type Properties map[string]string

// Page is structure holding HTML page metadata.
type Page struct {
	Lang       string     `json:"lang,omitempty"`
	Properties Properties `json:"properties,omitempty"`

	// TODO(mr) Support for prefixes
	prefixes map[string]string
}

// Title returns the page title based on defined priorities (html 5 > dc > og > twitter > title)
func (p Page) Title() string {
	propNames := []string{"dc:title", "og:title", "twitter:title", "title"}
	for _, name := range propNames {
		value := p.Properties[name]
		if value != "" {
			return value
		}
	}
	return ""
}

// ReadPage is used to extract metadata from an HTML page.
// It returns a Page struct for easy manipulation of those metadata.
func ReadPage(body io.Reader) (Page, error) {
	var p Page
	p.Properties = make(map[string]string)

	tokenizer := html.NewTokenizer(body)
Loop:
	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err != io.EOF {
				return p, err
			}
			break Loop

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "meta":
				meta := extract(token)
				if contains(knownProperties, meta.property) {
					p.Properties[meta.property] = meta.content
				}
			case "title":
				// The next token should be the page title
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					// Use page title but keep on searching an RDFa or Open Graph title, which is often more accurate
					p.Properties["title"] = tokenizer.Token().Data
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

	return p, nil
}

// ExtractRelMe
// TODO We also need to extract profiles from linked RDF cards.
func ExtractRelMe(body io.Reader) ([]string, error) {
	var urls []string

	tokenizer := html.NewTokenizer(body)
Loop:
	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err != io.EOF {
				return urls, err
			}
			break Loop

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "link", "a":
				relUrl, matched := matchAttr(token, "rel", "me", "href")
				if matched {
					urls = append(urls, relUrl)
				}
			}
		}
	}

	return urls, nil
}

// matchAttr returns the value of a given attributes, assuming we match the value of another one.
// For example, to extract href value on a rel link, you can call:
// value, matched := matchAttr(token, "rel", "me", "href")
func matchAttr(token html.Token, attrName, keywordToMatch, attrForValue string) (string, bool) {
	value := ""
	matched := false
	for _, attr := range token.Attr {
		switch attr.Key {
		case attrName:
			keyValue := attr.Val

			keyContent := strings.Split(keyValue, " ")
			if contains(keyContent, keywordToMatch) {
				matched = true
			}
		case attrForValue:
			value = attr.Val
		}
		if attr.Key == attrName {
		}
	}
	return value, matched
}

//============================================================================
// Properties extraction

var knownProperties = []string{
	// Dublin Core (HTML 5)
	"dc:title", "dc:creator",
	// Open Graph
	"og:title", "og:type", "og:url", "og:image",
	"og:description", "og:site_name",
	// Twitter
	"twitter:card", "twitter:site", "twitter:title",
	"twitter:image", "twitter:description",
	// Extra real world usage
	"description",
}

type meta struct {
	property string
	content  string
}

func extract(token html.Token) meta {
	var m meta

	for _, attr := range token.Attr {
		if attr.Key == "property" {
			m.property = attr.Val
		}
		// Twitter is incorrectly using name attribute to hold metadata
		// For details, see: https://www.ctrl.blog/entry/rdfa-socialmedia-metadata
		if m.property == "" && attr.Key == "name" {
			m.property = attr.Val
		}
		if attr.Key == "content" {
			m.content = attr.Val
		}
	}
	return m
}

//============================================================================
// Helper functions

func contains(array []string, str string) bool {
	for _, elt := range array {
		if elt == str {
			return true
		}
	}
	return false
}

//=============================================================================
// TODO

// References:
// - Contains example for XHTML and for setting metadata outside of HTML head
//   https://www.w3.org/MarkUp/2009/rdfa-for-html-authors

// TODO also extract og:image. e.g.:
// <meta property="og:image" content="https://gigaom.com/wp-content/uploads/sites/1/2011/01/sonosgroup-804x516.jpg" />

// TODO Add support for older Dublin Core syntax.
// See: https://www.slideshare.net/eduservfoundation/dublin-core-basic-syntax-tutorial
