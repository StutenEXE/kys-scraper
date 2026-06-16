package results

type Edition struct {
	ISBN13      string
	ISBN10      string
	Publisher   string
	PublishDate string
	PageCount   int
	Cover       string
	Dimensions  struct {
		Height    string
		Width     string
		Thickness string
	}
}

func (e Edition) ResultType() string { return "edition" }

func (e Edition) ToMap() map[string]any {
	return map[string]any{
		"isbn13":      e.ISBN13,
		"isbn10":      e.ISBN10,
		"publisher":   e.Publisher,
		"publishDate": e.PublishDate,
		"pageCount":   e.PageCount,
		"cover":       e.Cover,
		"dimensions": map[string]any{
			"height":    e.Dimensions.Height,
			"width":     e.Dimensions.Width,
			"thickness": e.Dimensions.Thickness,
		},
	}
}
