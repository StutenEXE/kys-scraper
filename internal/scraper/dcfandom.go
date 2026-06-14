package scraper

import "scrapers/internal/scraper/fandom"

func NewDCFandomScraper() *FandomScraper {
	return NewFandomScraper("dc.fandom.com", fandom.DefaultClassifiers()...)
}
