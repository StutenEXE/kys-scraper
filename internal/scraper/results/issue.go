package results

type Issue struct {
	Name         string
	Number       string
	ParutionDate string
	CoverDate    string
}

func (c Issue) ResultType() string { return "issue" }

func (c Issue) ToMap() map[string]string {
	return map[string]string{
		"name":         c.Name,
		"number":       c.Number,
		"parutionDate": c.ParutionDate,
		"coverDate":    c.CoverDate,
	}
}
