package semweb

import (
	"io"
	"sync"
)

// Extractor is a function that is used to define the process the crawler
// will execute.
type Extractor func(io.ReadCloser, Client) ([]string, []interface{})

type Processor interface {
	Process(io.Reader, Client) []string
}

// Crawler is a tool to crawl web URLs
type Crawler struct {
	client    Client
	queue     chan string
	wg        sync.WaitGroup
	processor Processor
	// TODO define max number of pages to crawl on a given domain to avoid being stuck
	// on a site generating a infinite number of random URLs.
}

func NewCrawler(proc Processor) *Crawler {
	queue := make(chan string)
	client := NewClient()
	crawler := Crawler{queue: queue, client: client, processor: proc}

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
// TODO: Skip url that were already checked.
func (c *Crawler) processURL(url string) {
	defer c.wg.Done()

	body, err := c.client.Get(url)
	if err != nil { // Cannot get URL
		return
	}

	// Pass body for page processing
	newURLs := c.processor.Process(body, c.client)
	body.Close()

	for _, u := range newURLs {
		c.enqueue(u)
	}
}
