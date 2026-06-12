package scraper

import "context"

// Scraper is implemented by every domain-specific scraper.
type Scraper interface {
	// Scrape fetches and extracts content from the given URL.
	Scrape(ctx context.Context, rawURL string) (ScrapeResult, error)
}
