package main

import (
	"fmt"
	"io"

	"github.com/processone/dpk/pkg/semweb"
)

//=============================================================================
// Test

type CrawlerProcessor struct {
	count       int
	profileUrls []string
	// TODO: handle profile validation (bidirectional profile validation)
}

func (cr *CrawlerProcessor) Process(body io.Reader, ctx semweb.Context) []string {
	cr.profileUrls = append(cr.profileUrls, ctx.Url)
	cr.count += 1

	urls, err := semweb.ExtractRelMe(ctx, body)
	if err != nil {
		fmt.Println("error:", err)
		return []string{}
	}

	if len(urls) == 0 {
		return []string{}
	}

	cleanUrls := make([]string, len(urls))
	for i, u := range urls {
		if targetUrl := ctx.Client.FollowRedirect(u); targetUrl != "" {
			cleanUrls[i] = targetUrl
		}
	}

	return cleanUrls
}

// Discover web profiles for a user, given a URL entrypoint
func main() {
	processor := CrawlerProcessor{}
	c := semweb.NewCrawler(&processor)
	c.Run("https://twitter.com/mickael")
	fmt.Println("Profiles:")
	for _, url := range processor.profileUrls {
		fmt.Println(url)
	}
}
