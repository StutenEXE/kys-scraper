package results

type Issue struct {
	Name         string
	Number       string
	ParutionDate string
	CoverDate    string
}

func (i Issue) ResultType() string { return "issue" }

func (i Issue) ToMap() map[string]any {
	return map[string]any{
		"name":         i.Name,
		"number":       i.Number,
		"parutionDate": i.ParutionDate,
		"coverDate":    i.CoverDate,
	}
}
