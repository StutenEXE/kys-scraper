package scraper

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"scrapers/internal/scraper/googlebooks"
	"scrapers/internal/scraper/results"
)

var (
	isbn13Re = regexp.MustCompile(`^97[89]\d{10}$`)
	isbn10Re = regexp.MustCompile(`^\d{9}[\dX]$`)
)

type ISBNScraper struct {
	client *googlebooks.Client
}

func NewISBNScraper() *ISBNScraper {
	return &ISBNScraper{
		client: googlebooks.NewClient(os.Getenv("GOOGLE_BOOKS_API_KEY")),
	}
}

func (s *ISBNScraper) Scrape(ctx context.Context, rawISBN string) (ScrapeResult, error) {
	isbn := normalizeISBN(rawISBN)

	isbnInfo, err := s.client.FetchByISBN(ctx, isbn)
	if err != nil {
		return ScrapeResult{}, fmt.Errorf("isbn scraper: %w", err)
	}

	return NewScrapeResult(results.Isbn{
		Book:    findBook(isbnInfo),
		Edition: findEdition(isbnInfo),
	}), nil

}

func normalizeISBN(s string) string {
	// Strip hyphens and spaces — 978-0-14-032872-1 → 9780140328721
	re := regexp.MustCompile(`[\s-]`)
	return re.ReplaceAllString(s, "")
}

func findBook(info *googlebooks.VolumeInfo) results.Book {
	return results.Book{
		Title:       strings.TrimSpace(info.Title + " " + info.Subtitle),
		Authors:     info.Authors,
		Cover:       info.ImageLinks.Thumbnail,
		Description: info.Description,
	}
}

func findEdition(info *googlebooks.VolumeInfo) results.Edition {
	ed := results.Edition{
		Publisher:   info.Publisher,
		PublishDate: info.PublishedDate,
		Cover:       info.ImageLinks.Thumbnail,
		PageCount:   info.PageCount,
	}
	for _, id := range info.IndustryIdentifiers {
		switch id.Type {
		case "ISBN_10":
			ed.ISBN10 = id.Identifier
		case "ISBN_13":
			ed.ISBN13 = id.Identifier
		}
	}
	return ed
}
