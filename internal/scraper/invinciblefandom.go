package scraper

import (
	"net/url"
	"regexp"
	"scrapers/internal/scraper/fandom"
	"scrapers/internal/scraper/results"
	"time"
)

func NewInvincibleFandomScraper() *FandomScraper {
	return NewFandomScraper("comic-invincible.fandom.com",
		invincibleIssueClassifier,
	)
}

var (
	invincibleIssuePattern = regexp.MustCompile(`/wiki/[^/]+_Issue_\d+$`)
)

func invincibleIssueClassifier(u *url.URL, data fandom.FandomData) (any, bool) {
	if !invincibleIssuePattern.MatchString(u.Path) {
		return nil, false
	}
	return results.Issue{
		Name:         invincibleFindIssueName(data),
		Number:       fandom.FindIssueNumber(data),
		ParutionDate: invincibleFindIssueParutionDate(data),
		CoverDate:    invincibleFindIssueParutionDate(data),
	}, true
}

func invincibleFindIssueName(data fandom.FandomData) string {
	text := data.Wikitext["body"]
	reIssueName := regexp.MustCompile(`== Appearing in "([^"]+)" ==`)
	matches := reIssueName.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func invincibleFindIssueParutionDate(data fandom.FandomData) string {
	datePub := data.Wikitext["datePublished"]
	t, _ := time.Parse("January, 2006", datePub)
	return t.Format("2006-01-02T15:04:05")
}

func invincibleFindIssueCoverDate(data fandom.FandomData) string {
	return invincibleFindIssueParutionDate(data)
}
