package fandom

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"scrapers/internal/scraper/helpers"
	"scrapers/internal/scraper/results"
)

var (
	issuePattern  = regexp.MustCompile(`/wiki/[^/]+_Vol_\d+_\d+$`)
	seriesPattern = regexp.MustCompile(`/wiki/[^/]+_Vol_\d+$`)
)

func IssueClassifier(u *url.URL, data FandomData) (any, bool) {
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

func IssueSeriesClassifier(u *url.URL, data FandomData) (any, bool) {
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

func DefaultClassifiers() []Classifier {
	return []Classifier{
		IssueClassifier,
		IssueSeriesClassifier,
	}
}

func WithDefaultClassifiers(extra ...Classifier) []Classifier {
	return append(DefaultClassifiers(), extra...)
}

func findIssueName(data FandomData) string {
	name := data.Fields["StoryTitle1"]
	if name == "" {
		// Build name from issue serie and issue number
		serie := data.Fields["Title"]
		number := findIssueNumber(data)
		if serie != "" && number != "" {
			name = fmt.Sprintf("%s #%s", serie, number)
		}
	}
	return name
}

func findIssueNumber(data FandomData) string {
	number := data.Fields["Issue"]
	if number == "" {
		// Ex : Ultimate Spider-Man Vol 1 [1]
		splitted := strings.Split(data.Title, " ")
		number = splitted[len(splitted)-1]
	}
	return number
}

func findIssueParutionDate(data FandomData) string {
	releaseDate := data.Fields["ReleaseDate"]
	if releaseDate != "" {
		// Release should be formatted as : September 7, 2000
		t, _ := time.Parse("January 2, 2006", releaseDate)
		return t.Format("2006-01-02T15:04:05")
	}
	day := data.Fields["Pubday"]
	if day == "" {
		day = data.Fields["Day"]
	}

	month := data.Fields["Pubmonth"]
	pubMonthPresent := month != ""
	if !pubMonthPresent {
		month = data.Fields["Month"]
		// Sometimes, Month can be written as January 2, and hold the date
		splittedMonth := strings.Split(month, " ")
		if len(splittedMonth) > 1 {
			month = splittedMonth[0]
			day = splittedMonth[1]
		}
	}

	year := data.Fields["Year"]
	t := helpers.ParseToDate(day, month, year)
	if !pubMonthPresent {
		t = t.AddDate(0, -2, 0)
	}
	return t.Format("2006-01-02T15:04:05")
}

func findIssueCoverDate(data FandomData) string {
	month := data.Fields["Month"]
	// Sometimes, Month can be written as January 2, and hold the date
	splittedMonth := strings.Split(month, " ")
	if len(splittedMonth) > 1 {
		month = splittedMonth[0]
	}
	year := data.Fields["Year"]
	t := helpers.ParseToDate("1", month, year)
	return t.Format("2006-01-02T15:04:05")
}

func findIssueSerieName(data FandomData) string {
	re := regexp.MustCompile(`(?i)\s+Vol\s+(\d+)$`)
	if re.MatchString(data.Title) {
		return re.ReplaceAllString(data.Title, " (Volume $1)")
	}
	return data.Title
}

func findIssueSerieDescription(data FandomData) string {
	history := data.Fields["History"]
	splitted := strings.Split(history, "\n\n")
	splitted = strings.Split(splitted[0], "==History==")
	reLinkApostrophes := regexp.MustCompile(`'{2,}`)
	text := reLinkApostrophes.ReplaceAllString(splitted[0], "")
	return strings.TrimSpace(text)
}

func findIssueSerieStartDate(data FandomData) string {
	month := data.Fields["StartMonth"]
	year := data.Fields["StartYear"]
	t := helpers.ParseToDate("1", month, year)
	return t.Format("2006-01-02T15:04:05")
}

func findIssueSerieEndDate(data FandomData) string {
	month := data.Fields["EndMonth"]
	year := data.Fields["EndYear"]
	if month == "" && year == "" {
		return ""
	}
	t := helpers.ParseToDate("1", month, year)
	return t.Format("2006-01-02T15:04:05")
}
