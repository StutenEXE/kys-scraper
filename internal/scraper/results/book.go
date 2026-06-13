package results

type Book struct {
	Title       string
	Authors     []string
	Cover       string
	Description string
}

func (b Book) ResultType() string { return "book" }

func (b Book) ToMap() map[string]any {
	return map[string]any{
		"title":       b.Title,
		"authors":     b.Authors,
		"cover":       b.Cover,
		"description": b.Description,
	}
}
