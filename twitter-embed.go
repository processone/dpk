package dpk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/processone/dpk/pkg/semweb"
)

func twitterEmbed(l Link, dataDir string) string {
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
	html2 := enrichHTML(safeHTML, dataDir)

	return "\n" + html2
}

func enrichHTML(fragment string, dataDir string) string {
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
		if newNode := walkNode(n, dataDir); newNode != nil {
			n.Parent.InsertBefore(newNode, n)
			n.Parent.RemoveChild(n)
		}
	}

	// ======================================
	// Render all elements from fragment.
	buf := new(bytes.Buffer)
	for _, n := range doc {
		if err := html.Render(buf, n); err != nil {
			// Return unmodified fragment on error
			fmt.Println("html render error:", err)
			return fragment
		}
	}
	// Return modified fragment
	return buf.String()
}

func walkNode(n *html.Node, dataDir string) *html.Node {
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
			if strings.HasPrefix(c.Data, "pic.twitter.com") { // Content rewrite
				// Download image and replace link with inline image:
				if img, err := DownloadImage(c.Data, dataDir); err == nil {
					var newAttr []html.Attribute
					imgSrc := append(newAttr, html.Attribute{"", "src", img})
					newNode := html.Node{
						Type:     html.ElementNode,
						DataAtom: atom.Img,
						Data:     "img",
						Attr:     imgSrc,
					}
					return &newNode
				}
			} else { // Link rewrite
				l.AnchorText = c.Data
				l = l.Resolve()
				if needRewrite(c.Data) {
					c.Data = l.URLTitle
				}

				// TODO: keep other attributes, like class
				var newAttr []html.Attribute
				n.Attr = append(newAttr, html.Attribute{"", "href", l.URL})
				return nil
			}
		}

		// Otherwise, resolve only the URL and keep the children untouched
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if newNode := walkNode(c, dataDir); newNode != nil {
				parent := c.Parent
				if c.NextSibling != nil {
					parent.RemoveChild(c)
					parent.InsertBefore(newNode, c.NextSibling)
				} else {
					parent.RemoveChild(c)
					parent.AppendChild(newNode)
				}
				c = newNode
			}
		}

		l = l.Resolve()
		n.Attr = []html.Attribute{{Key: "href", Val: l.URL}}
		return nil
	}

	// Just walk down the node structure if this is not a link node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if newNode := walkNode(c, dataDir); newNode != nil {
			parent := c.Parent
			if c.NextSibling != nil {
				parent.RemoveChild(c)
				parent.InsertBefore(newNode, c.NextSibling)
			} else {
				parent.RemoveChild(c)
				parent.AppendChild(newNode)
			}
			c = newNode
		}
	}
	return nil
}

func needRewrite(urlText string) bool {
	if strings.HasPrefix(urlText, "http://") {
		return true
	}
	if strings.HasPrefix(urlText, "https://") {
		return true
	}
	return false
}

// DownloadImage downloads a Twitter image based on its pic.twitter.com URL.
// TODO(mr): Move twitter image manipulation into its own package.
func DownloadImage(imageRef, dataDir string) (string, error) {
	imageURL := imageRef
	if !strings.HasPrefix(imageRef, "http") {
		imageURL = "https://" + imageRef
	}

	targetUrl, err := GetImageURL(imageURL)
	if err != nil {
		return imageURL, err
	}

	resp, err := httpClient().Get(targetUrl)
	if err != nil {
		return imageURL, err
	}
	defer resp.Body.Close()

	filename := filepath.Join(dataDir, filepath.Base(targetUrl))
	file, err := os.Create(filename)
	if err != nil {
		return imageURL, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return imageURL, err
	}
	return filepath.Base(targetUrl), nil
}

func GetImageURL(imageUrl string) (string, error) {
	// Discover image URL
	client := semweb.NewClient()
	body, _, err := client.Get(imageUrl)
	if err != nil {
		return imageUrl, err
	}
	defer body.Close()

	page, err := semweb.ReadPage(body)
	if err != nil {
		return imageUrl, err
	}

	img := page.Properties["og:image"]
	return img, nil
}
