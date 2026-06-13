// internal/scraper/results/comic_issue.go
package results

type IssueSerie struct {
	Name        string
	StartDate   string
	EndDate     string
	Description string
}

func (is IssueSerie) ResultType() string { return "issueserie" }

func (is IssueSerie) ToMap() map[string]any {
	return map[string]any{
		"name":        is.Name,
		"startDate":   is.StartDate,
		"endDate":     is.EndDate,
		"description": is.Description,
	}
}
