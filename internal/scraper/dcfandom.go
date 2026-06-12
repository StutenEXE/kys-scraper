package scraper

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"scrapers/internal/scraper/fandom"
	"scrapers/internal/scraper/helpers"
	"scrapers/internal/scraper/results"
)

// classifier inspects a URL and loaded document, and returns a result if it
// recognises the page. Returns (zero, false) to pass to the next classifier.
type classifier func(u *url.URL, data fandom.FandomData) (TypedResult, bool)

type DcFandomScraper struct {
	client      *fandom.Client
	classifiers []classifier
}

func NewDCFandomScraper() *DcFandomScraper {
	return &DcFandomScraper{
		client: fandom.NewClient("dc.fandom.com"),
		classifiers: []classifier{
			classifyIssue,       // e.g. /wiki/Batman_Vol_1_666
			classifyIssueSeries, // e.g. /wiki/Batman_Vol_1
		},
	}
}

var (
	issuePattern  = regexp.MustCompile(`/wiki/[^/]+_Vol_\d+_\d+$`)
	seriesPattern = regexp.MustCompile(`/wiki/[^/]+_Vol_\d+$`)
)

func (s *DcFandomScraper) Scrape(ctx context.Context, rawURL string) (ScrapeResult, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ScrapeResult{}, fmt.Errorf("invalid url: %w", err)
	}

	// Extract page title from path, e.g. /wiki/Batman_Vol_1_666 → Batman_Vol_1_666
	pageTitle := strings.TrimPrefix(u.Path, "/wiki/")

	page, err := s.client.FetchPage(ctx, pageTitle)
	if err != nil {
		return ScrapeResult{}, err
	}

	wikitext := page.Parse.Wikitext["*"]
	fields := fandom.ParseInfobox(wikitext)

	data := fandom.FandomData{
		Title:  page.Parse.Title,
		Fields: fields,
	}

	for _, classify := range s.classifiers {
		if typed, ok := classify(u, data); ok {
			return NewScrapeResult(typed), nil
		}
	}

	return ScrapeResult{}, fmt.Errorf("dcfandom: no classifier matched %s", rawURL)
}

// classifyIssue matches /wiki/Batman_Vol_1_666
func classifyIssue(u *url.URL, data fandom.FandomData) (TypedResult, bool) {
	if !issuePattern.MatchString(u.Path) {
		return nil, false
	}
	return results.Issue{
		Name:         findIssueName(data),
		Number:       findIssueNumber(data),
		ParutionDate: findIssueParutionDate(data),
		CoverDate:    findIssueCoverDate(data),
	}, true
}

// classifyIssueSeries matches /wiki/Batman_Vol_1
func classifyIssueSeries(u *url.URL, data fandom.FandomData) (TypedResult, bool) {
	if !seriesPattern.MatchString(u.Path) {
		return nil, false
	}
	return results.IssueSerie{
		Name:        findIssueSerieName(data),
		Description: findIssueSerieDescription(data),
		StartDate:   findIssueSerieStartDate(data),
		EndDate:     findIssueSerieEndDate(data),
	}, true
}

func findIssueName(data fandom.FandomData) string {
	return data.Fields["StoryTitle1"]
}

func findIssueNumber(data fandom.FandomData) string {
	return data.Fields["Issue"]
}

func findIssueParutionDate(data fandom.FandomData) string {
	// Always a number or is not present
	day := data.Fields["Day"]
	// A number of a string, should always be present, generally 2 months ahead (cover date)
	month := data.Fields["Month"]
	// A number, should always be present
	year := data.Fields["Year"]
	t := helpers.ParseToDate(day, month, year)
	// Go back 2 months
	t = t.AddDate(0, -2, 0)
	return t.Format("2006-01-02T15:04:05")
}

func findIssueCoverDate(data fandom.FandomData) string {
	// A number of a string, should always be present
	month := data.Fields["Month"]
	// A number, should always be present
	year := data.Fields["Year"]
	t := helpers.ParseToDate("1", month, year)

	return t.Format("2006-01-02T15:04:05")
}

func findIssueSerieName(data fandom.FandomData) string {
	return data.Title
}

func findIssueSerieDescription(data fandom.FandomData) string {
	splitted := strings.Split(data.Fields["History"], "\n\n")
	return splitted[0]
}

func findIssueSerieStartDate(data fandom.FandomData) string {
	// A number of a string, should always be present, generally 2 months ahead (cover date)
	month := data.Fields["StartMonth"]
	// A number, should always be present
	year := data.Fields["StartYear"]
	t := helpers.ParseToDate("1", month, year)
	return t.Format("2006-01-02T15:04:05")
}

func findIssueSerieEndDate(data fandom.FandomData) string {
	// A number of a string, should always be present
	month := data.Fields["EndMonth"]
	// A number, should always be present
	year := data.Fields["EndYear"]
	t := helpers.ParseToDate("1", month, year)
	return t.Format("2006-01-02T15:04:05")
}
