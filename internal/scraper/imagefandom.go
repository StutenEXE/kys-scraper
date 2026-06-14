package scraper

import (
	"net/url"
	"scrapers/internal/scraper/fandom"
	"scrapers/internal/scraper/helpers"
	"scrapers/internal/scraper/results"
)

func NewImageFandomScraper() *FandomScraper {
	return NewFandomScraper("imagecomics.fandom.com",
		fandom.IssueSeriesClassifier,
		imageIssueClassifier,
	)
}

func imageIssueClassifier(u *url.URL, data fandom.FandomData) (any, bool) {
	if !fandom.IssuePattern.MatchString(u.Path) {
		return nil, false
	}
	return results.Issue{
		Name:         fandom.FindIssueName(data),
		Number:       fandom.FindIssueNumber(data),
		ParutionDate: imageFindIssueParutionDate(data),
		CoverDate:    fandom.FindIssueCoverDate(data),
	}, true
}

func imageFindIssueParutionDate(data fandom.FandomData) string {
	day := data.Fields["Day"]
	month := data.Fields["Month"]
	year := data.Fields["Year"]
	t := helpers.ParseToDate(day, month, year)
	return t.Format("2006-01-02T15:04:05")
}
