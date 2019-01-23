package dpk

import (
	"fmt"

	"github.com/processone/dpk/pkg/semweb"
)

type Link struct {
	URL string
	// Display text is the text of the original link
	AnchorText string
	// whereas URLTitle is the title discovered on the actual URL
	URLTitle string

	resolved bool
}

func (l Link) Resolve() Link {
	// Link was already processed:
	if l.resolved {
		return l
	}

	fmt.Println("Processing link:", l.URL)
	l.resolved = true

	client := semweb.NewClient()
	body, err := client.Get(l.URL) // TODO: Should return a Result having a body (readcloser and an actual URL). What about canonical URLs
	if err != nil {
		return l
	}
	defer body.Close()

	// TODO: We need the client to return the final URL of the Get:
	// l.URL = finalUrl
	page, err := semweb.ReadPage(body)
	if err == nil {
		l.URLTitle = page.Title()
	}

	return l
}

// TODO: Add that method to a renderer. We should have at least a Markdown and an HTML renderer.
func (l Link) Markdown() string {
	text := l.AnchorText
	if l.URLTitle != "" {
		text = l.URLTitle
	}

	if len(text) > 50 {
		text = text[:50] + "…"
	}
	return fmt.Sprintf("[%s](%s)", text, l.URL)
}
