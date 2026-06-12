package registry

import (
	"scrapers/internal/dispatcher"
	"scrapers/internal/scraper"
)

func All() []dispatcher.Registration {
	return []dispatcher.Registration{
		{
			Match:   func(host string) bool { return host == "dc.fandom.com" },
			Factory: func() scraper.Scraper { return scraper.NewDCFandomScraper() },
		},
	}
}
