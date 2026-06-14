package scraper

import (
	"context"
	"net/url"
	"scrapers/internal/scraper/fandom"
)

// classifier inspects a URL and loaded document, and returns a result if it
// recognises the page. Returns (zero, false) to pass to the next classifier.
type classifier func(u *url.URL, data fandom.FandomData) (TypedResult, bool)

// Scraper is implemented by every domain-specific scraper.
type Scraper interface {
	// Scrape fetches and extracts content from the given URL.
	Scrape(ctx context.Context, rawURL string) (ScrapeResult, error)
}
