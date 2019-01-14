package semweb

import (
	"io"
	"sync"
)

type Context struct {
	Client Client
	Url    string
}

type Processor interface {
	Process(io.Reader, Context) []string
}

// Crawler is a tool to crawl web URLs
type Crawler struct {
	client    Client
	queue     chan string
	wg        sync.WaitGroup
	processor Processor
	visited   map[string]bool
	// TODO define max number of pages to crawl on a given domain to avoid being stuck
	// on a site generating a infinite number of random URLs.
}

func NewCrawler(proc Processor) *Crawler {
	queue := make(chan string)
	// TODO: Add ability to customize HTTP client => Use option parameter
	client := NewClient()
	crawler := Crawler{queue: queue, client: client,
		processor: proc, visited: make(map[string]bool)}

	// Main processing loop
	go func() {
		for url := range queue {
			crawler.processURL(url)
		}
	}()

	return &crawler
}

func (c *Crawler) Run(url string) {
	c.wg.Add(1)
	c.queue <- url
	c.wg.Wait()
}

// Enqueue start the crawling from a given URL or add a new discovered URL to the
// queue.
func (c *Crawler) enqueue(url string) {
	c.wg.Add(1)
	// We enqueue asynchronously to avoid locking
	go func() {
		c.queue <- url
	}()
}

// processURL retrieves a give URL and pass it to the features extractor.
// TODO:
//   - Skip urls that were already checked.
//   - Store url and their canonical URLs ? check how to best handle canonical url
func (c *Crawler) processURL(url string) {
	defer c.wg.Done()

	c.visited[url] = true
	body, err := c.client.Get(url)
	if err != nil { // Cannot get URL
		return
	}

	// Pass body for page processing and context for proper page analysis, relative link resolution, etc.
	context := Context{Client: c.client, Url: url}
	newURLs := c.processor.Process(body, context)
	body.Close()

	for _, u := range newURLs {
		if !c.visited[u] {
			c.enqueue(u)
		}
	}
}
