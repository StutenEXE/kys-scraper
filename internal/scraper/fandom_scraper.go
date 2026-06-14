package scraper

import (
	"context"
	"fmt"

	"scrapers/internal/scraper/fandom"
)

type FandomScraper struct {
	host   string
	client *fandom.HostScraper
}

func NewFandomScraper(host string, classifiers ...fandom.Classifier) *FandomScraper {
	return &FandomScraper{
		host:   host,
		client: fandom.NewHostScraper(host, classifiers...),
	}
}

func (s *FandomScraper) Scrape(ctx context.Context, rawURL string) (ScrapeResult, error) {
	result, err := s.client.Scrape(ctx, rawURL)
	if err != nil {
		return ScrapeResult{}, err
	}

	typed, ok := result.(TypedResult)
	if !ok {
		return ScrapeResult{}, fmt.Errorf("%s: unsupported result type", s.host)
	}

	return NewScrapeResult(typed), nil
}
