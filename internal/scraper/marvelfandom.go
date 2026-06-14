package scraper

import "scrapers/internal/scraper/fandom"

func NewMarvelFandomScraper() *FandomScraper {
	return NewFandomScraper("marvel.fandom.com", fandom.DefaultClassifiers()...)
}
