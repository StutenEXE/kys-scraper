package results

type Edition struct {
	ISBN13      string
	ISBN10      string
	Publisher   string
	PublishDate string
	PageCount   int
	Cover       string
}

func (e Edition) ResultType() string { return "edition" }

func (e Edition) ToMap() map[string]any {
	return map[string]any{
		"isbn13":      e.ISBN13,
		"isbn10":      e.ISBN13,
		"publisher":   e.Publisher,
		"publishDate": e.PublishDate,
		"pageCount":   e.PageCount,
		"cover":       e.Cover,
	}
}
