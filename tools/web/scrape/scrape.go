package webscrapetool

import (
	"time"

	"github.com/gocolly/colly/v2"
)

type ScrapeInput struct {
	Urls []string `json:"urls" jsonschema:"urls to scrape"`
}

type ScrapedPage struct {
	Url     string
	Title   string
	Content string
	Err     error
}

func getCollyCollector() (*colly.Collector, error) {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0"),
	)
	c.SetRequestTimeout(15 * time.Second)
	if err := c.Limit(&colly.LimitRule{
		Parallelism: 3,
	}); err != nil {
		return nil, err
	}

	return c, nil
}
