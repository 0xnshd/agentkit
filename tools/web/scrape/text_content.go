package webscrapetool

import (
	"bytes"
	"strings"

	"codeberg.org/readeck/go-readability/v2"
	"github.com/gocolly/colly/v2"
)

func collyScrape(urls []string, results chan<- ScrapedPage) {

	c, err := getCollyCollector()
	if err != nil {
		for _, url := range urls {
			results <- ScrapedPage{Url: url, Err: err}
		}
		return
	}

	c.OnResponse(func(r *colly.Response) {
		article, err := readability.FromReader(
			bytes.NewReader(r.Body),
			r.Request.URL,
		)

		sp := ScrapedPage{Url: r.Request.URL.String()}
		if err != nil {
			sp.Err = err
		} else {
			sp.Title = article.Title()

			var b strings.Builder
			if err := article.RenderText(&b); err != nil {
				sp.Err = err
			} else {
				sp.Content = b.String()
			}
		}
		results <- sp
	})

	c.OnError(func(r *colly.Response, err error) {
		results <- ScrapedPage{Url: r.Request.URL.String(), Err: err}
	})

	for _, url := range urls {
		if err := c.Visit(url); err != nil {
			results <- ScrapedPage{Url: url, Err: err}
		}
	}

	c.Wait()
}

func ScrapeTextContent(urls []string) <-chan ScrapedPage {
	results := make(chan ScrapedPage, len(urls))

	go func() {
		defer close(results)
		collyScrape(urls, results)
	}()

	return results
}
