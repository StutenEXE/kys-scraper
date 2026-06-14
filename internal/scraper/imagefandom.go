package scraper

import "scrapers/internal/scraper/fandom"

func NewImageFandomScraper() *FandomScraper {
	return NewFandomScraper("imagecomics.fandom.com", fandom.DefaultClassifiers()...)
}
