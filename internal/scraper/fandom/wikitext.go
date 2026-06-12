package fandom

import (
	"regexp"
	"strings"
)

var infoboxFieldRe = regexp.MustCompile(`\|\s*(\w+)\s*=\s*([^\|}\n]+)`)

// ParseInfobox extracts key/value pairs from a wikitext infobox.
func ParseInfobox(wikitext string) map[string]string {
	fields := make(map[string]string)
	matches := infoboxFieldRe.FindAllStringSubmatch(wikitext, -1)
	for _, m := range matches {
		key := strings.TrimSpace(m[1])
		val := cleanWikitext(strings.TrimSpace(m[2]))
		fields[key] = val
	}
	return fields
}

// cleanWikitext strips wiki markup like [[links]], {{templates}}, and HTML tags.
func cleanWikitext(s string) string {
	// [[Display Text|Link]] or [[Link]] → Display Text or Link
	s = regexp.MustCompile(`(?s)\[\[(?:[^\]|]*\|)?([^\]|]+)\]\]`).ReplaceAllString(s, "$1")
	// {{template}} → ""
	s = regexp.MustCompile(`\{\{[^}]*\}\}`).ReplaceAllString(s, "")
	// <ref>...</ref> → ""
	s = regexp.MustCompile(`<ref[^>]*>.*?</ref>`).ReplaceAllString(s, "")
	// remaining HTML tags
	s = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(s, "")

	return strings.TrimSpace(s)
}

// SplitList splits a wikitext list value into individual items.
// Handles both comma-separated and <br/>-separated lists.
func SplitList(s string) []string {
	s = regexp.MustCompile(`<br\s*/?>`).ReplaceAllString(s, ",")
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
