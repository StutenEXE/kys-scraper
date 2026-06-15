package fandom

import (
	"html"
	"regexp"
	"strings"
)

// extractRawPairs scans wikitext for every top-level {{...}} template block
// and extracts all named (key=value) parameters, merging the results from
// every block into a single map. Values are returned exactly as found
// (raw) — a separate sanitize step is expected to clean them up.
//
// Example input:
//
//	{{DISPLAYTITLE:Invincible Issue #144}}
//	{{Book|title=Invincible Issue #144|author=Robert Kirkman|...}}
//	{{#SQuote:|text=...|speaker=[[Mark Grayson|Invincible]]|...}}
//
// Output contains "title", "author", ... from {{Book}} AND
// "text", "speaker", ... from {{#SQuote:}} all in the same map.
func extractRawPairs(wikitext string) map[string]string {
	pairs := make(map[string]string)

	i := 0
	for i < len(wikitext)-1 {
		if wikitext[i] == '{' && wikitext[i+1] == '{' {
			end := findTemplateEnd(wikitext, i)
			content := wikitext[i+2 : end-2]
			extractTemplatePairs(content, pairs)
			i = end
			continue
		}
		i++
	}

	return pairs
}

// findTemplateEnd returns the index just past the closing "}}" that
// matches the "{{" starting at position start, accounting for any nested
// {{...}} templates within.
func findTemplateEnd(s string, start int) int {
	depth := 1
	i := start + 2
	for i < len(s)-1 {
		switch {
		case s[i] == '{' && s[i+1] == '{':
			depth++
			i += 2
		case s[i] == '}' && s[i+1] == '}':
			depth--
			i += 2
			if depth == 0 {
				return i
			}
		default:
			i++
		}
	}
	return len(s)
}

// extractTemplatePairs parses the content of a single template (without the
// surrounding {{ }}) and merges its named parameters into pairs.
func extractTemplatePairs(content string, pairs map[string]string) {
	segments := splitTopLevel(content, '|')
	if len(segments) <= 1 {
		return // no parameters, just a template name (e.g. DISPLAYTITLE)
	}

	// segments[0] is the template name (e.g. "Book", "#SQuote:");
	// the rest are parameters.
	for _, seg := range segments[1:] {
		eqIdx := strings.Index(seg, "=")
		if eqIdx == -1 {
			continue // positional/unnamed parameter, skip
		}

		key := strings.TrimSpace(seg[:eqIdx])
		if key == "" {
			continue
		}

		pairs[key] = seg[eqIdx+1:]
	}
}

// splitTopLevel splits s on sep, but ignores occurrences of sep nested
// inside {{...}}, [[...]], or <gallery>...</gallery> so that pipes
// belonging to inner templates, links, or galleries don't break apart
// outer parameters.
func splitTopLevel(s string, sep byte) []string {
	var segments []string
	var current strings.Builder
	depth := 0

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch {
		case c == '{' && i+1 < len(s) && s[i+1] == '{':
			depth++
			current.WriteByte(c)
			current.WriteByte(s[i+1])
			i++
			continue
		case c == '}' && i+1 < len(s) && s[i+1] == '}':
			if depth > 0 {
				depth--
			}
			current.WriteByte(c)
			current.WriteByte(s[i+1])
			i++
			continue
		case c == '[' && i+1 < len(s) && s[i+1] == '[':
			depth++
			current.WriteByte(c)
			current.WriteByte(s[i+1])
			i++
			continue
		case c == ']' && i+1 < len(s) && s[i+1] == ']':
			if depth > 0 {
				depth--
			}
			current.WriteByte(c)
			current.WriteByte(s[i+1])
			i++
			continue
		case strings.HasPrefix(strings.ToLower(s[i:]), "<gallery"):
			depth++
			for i < len(s) && s[i] != '>' {
				current.WriteByte(s[i])
				i++
			}
			if i < len(s) {
				current.WriteByte(s[i]) // '>'
			}
			continue
		case strings.HasPrefix(strings.ToLower(s[i:]), "</gallery>"):
			if depth > 0 {
				depth--
			}
			current.WriteString(s[i : i+10])
			i += 9
			continue
		}

		if c == sep && depth == 0 {
			segments = append(segments, current.String())
			current.Reset()
			continue
		}

		current.WriteByte(c)
	}

	segments = append(segments, current.String())
	return segments
}

// extractBody strips any leading top-level {{...}} template blocks from
// wikitext and returns the remaining content — i.e. the article body
// (headings, character lists, plot synopsis, categories, etc.) with
// leading/trailing whitespace trimmed. This is the "outer text" that sits
// outside the infobox/quote templates.
func extractBody(wikitext string) string {
	i := 0
	for i < len(wikitext)-1 {
		if wikitext[i] == '{' && wikitext[i+1] == '{' {
			i = findTemplateEnd(wikitext, i)
			continue
		}
		break
	}

	return strings.TrimSpace(wikitext[i:])
}

// ParseWikitext extracts key/value pairs from a wikitext infobox.
func ParseWikitext(wikitext string) map[string]string {
	fields := make(map[string]string)
	raw := extractRawPairs(wikitext)
	raw["body"] = extractBody(wikitext)
	fields = sanitizeAll(raw)
	return fields
}

var (
	// HTML comments: <!-- ... -->
	reHTMLComment = regexp.MustCompile(`<!--.*?-->`)
	// Wiki templates: {{ ... }}
	reTemplate = regexp.MustCompile(`\{\{[^}]*\}\}`)
	// Wiki links: [[Target|Label]] → Label, or [[Target]] → Target
	reWikiLink = regexp.MustCompile(`\[\[(?:[^|\]]*\|)?([^\]]+)\]\]`)
	// HTML tags
	reHTMLTag = regexp.MustCompile(`<[^>]+>`)
	// Collapse multiple spaces/newlines
	reWhitespace = regexp.MustCompile(`\s+`)
	// Remove remaining open/closes double brackets
	reDoubleBrackets = regexp.MustCompile(`[{]{2}|[}]{2}`)
)

func sanitize(raw string) string {
	s := raw
	// Decode HTML entities first (&amp, &lt, ...)
	s = html.UnescapeString(s)
	// Remove HTML comments
	s = reHTMLComment.ReplaceAllString(s, "")
	// Remove wiki templates {{...}}
	s = reTemplate.ReplaceAllString(s, "")
	// Resolve wiki links [[Target|Label]] → Label
	s = reWikiLink.ReplaceAllString(s, "$1")
	// Strip HTML tags
	s = reHTMLTag.ReplaceAllString(s, " ")
	// Trim and collapse whitespace
	s = reWhitespace.ReplaceAllString(s, " ")
	// Remove double brackets
	s = reDoubleBrackets.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

func sanitizeAll(raw map[string]string) map[string]string {
	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = sanitize(v)
	}
	return result
}
