package results

type Isbn struct {
	Book    Book
	Edition Edition
}

func (isbn Isbn) ResultType() string { return "isbn" }

func (isbn Isbn) ToMap() map[string]any {
	return map[string]any{
		"book":    isbn.Book.ToMap(),
		"edition": isbn.Edition.ToMap(),
	}
}
