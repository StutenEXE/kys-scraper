package scraper

// TypedResult is implemented by every concrete result struct
// It enforces that each type declares its own ResultType string
// and knows how to flatten itself into the wire format
type TypedResult interface {
	ResultType() string
	ToMap() map[string]string
}

// ScrapeResult is the envelope sent over the wire.
// Build it only via NewScrapeResult, never construct manually.
type ScrapeResult struct {
	ResultType string            `json:"resultType"`
	Result     map[string]string `json:"result"`
}

func NewScrapeResult(r TypedResult) ScrapeResult {
	return ScrapeResult{
		ResultType: r.ResultType(),
		Result:     r.ToMap(),
	}
}
