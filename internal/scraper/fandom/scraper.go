package fandom

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type Classifier func(u *url.URL, data FandomData) (any, bool)

type Scraper struct {
	client *Client
}

type HostScraper struct {
	*Scraper
	classifiers []Classifier
}

func NewScraper(wikiHost string) *Scraper {
	return &Scraper{
		client: NewClient(wikiHost),
	}
}

func NewHostScraper(wikiHost string, classifiers ...Classifier) *HostScraper {
	return NewScraper(wikiHost).Wrap(classifiers...)
}

func (s *Scraper) Wrap(classifiers ...Classifier) *HostScraper {
	if len(classifiers) == 0 {
		classifiers = DefaultClassifiers()
	}
	return &HostScraper{
		Scraper:     s,
		classifiers: classifiers,
	}
}

func (s *Scraper) Load(ctx context.Context, rawURL string) (FandomData, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return FandomData{}, fmt.Errorf("invalid url: %w", err)
	}

	pageTitle := strings.TrimPrefix(u.Path, "/wiki/")
	if pageTitle == "" {
		return FandomData{}, fmt.Errorf("invalid wiki path %s", rawURL)
	}

	page, err := s.client.FetchPage(ctx, pageTitle)
	if err != nil {
		return FandomData{}, err
	}

	wikitext, ok := page.Parse.Wikitext["*"]
	if !ok {
		return FandomData{}, fmt.Errorf("missing wikitext for %s", rawURL)
	}
	wtFields := ParseWikitext(wikitext)

	return FandomData{
		Title:    page.Parse.Title,
		Wikitext: wtFields,
	}, nil
}

func (s *Scraper) Scrape(ctx context.Context, rawURL string, classifiers ...Classifier) (any, error) {
	data, err := s.Load(ctx, rawURL)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	for _, classify := range classifiers {
		if result, ok := classify(u, data); ok {
			return result, nil
		}
	}

	return nil, fmt.Errorf("no classifier matched %s", rawURL)
}

func (h *HostScraper) Scrape(ctx context.Context, rawURL string) (any, error) {
	return h.Scraper.Scrape(ctx, rawURL, h.classifiers...)
}
