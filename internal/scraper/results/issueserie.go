// internal/scraper/results/comic_issue.go
package results

type IssueSerie struct {
	Name        string
	StartDate   string
	EndDate     string
	Description string
}

func (c IssueSerie) ResultType() string { return "issueserie" }

func (c IssueSerie) ToMap() map[string]string {
	return map[string]string{
		"name":        c.Name,
		"startDate":   c.StartDate,
		"endDate":     c.EndDate,
		"description": c.Description,
	}
}
