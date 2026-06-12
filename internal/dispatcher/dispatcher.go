package dispatcher

import (
	"fmt"
	"net/url"
	"scrapers/internal/scraper"
)

// Registration pairs a host matcher with a scraper factory
type Registration struct {
	// Match returns true when this scraper should handle the given host
	Match   func(host string) bool
	Factory func() scraper.Scraper
}

// Dispatcher holds an ordered list of registrations.
type Dispatcher struct {
	registrations []Registration
}

func New(regs []Registration) *Dispatcher {
	return &Dispatcher{registrations: regs}
}

// For returns the first scraper whose Match function accepts the URL's host
// Falls back to the last registration (expected to be the generic scraper)
func (d *Dispatcher) For(rawURL string) (scraper.Scraper, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	for _, reg := range d.registrations {
		if reg.Match(u.Hostname()) {
			return reg.Factory(), nil
		}
	}
	return nil, fmt.Errorf("no scraper found for %s", u.Hostname())
}
