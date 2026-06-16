package scraper

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	// For dimensions, we tranform to cm if needed
	ed.Dimensions = results.EditionDimensions{
		Height:    parseDimensionToCm(info.Dimensions.Height),
		Width:     parseDimensionToCm(info.Dimensions.Width),
		Thickness: parseDimensionToCm(info.Dimensions.Thickness),
	}
	return ed
}

// parseDimensionToCm takes a string like "26.50 cm" or "10.43 in"
// and returns it normalized to cm, e.g. "26.50".
func parseDimensionToCm(dim string) float64 {
	dim = strings.TrimSpace(dim)
	if dim == "" {
		return 0
	}

	parts := strings.Fields(dim) // splits on whitespace, handles multiple spaces too
	if len(parts) != 2 {
		// Unexpected format, return as-is rather than guessing
		return 0
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	switch strings.ToLower(parts[1]) {
	case "in", "in.", "inch", "inches":
		value *= 2.54
	case "cm":
		// already correct unit, no-op
	default:
		// unknown unit, return original untouched
		return 0
	}

	return value
}
