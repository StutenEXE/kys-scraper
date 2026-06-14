package fandom

import (
	"html"
	"regexp"
	"strings"
)

func extractRawPairs(wikitext string) map[string]string {
	result := make(map[string]string)

	lines := strings.Split(wikitext, "\n")
	var currentKey string
	var currentVal strings.Builder

	// flush stores the current key-value pair and resets state.
	flush := func() {
		if currentKey != "" {
			result[currentKey] =
				strings.TrimSpace(currentVal.String())
		}
		currentKey = ""
		currentVal.Reset()
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			// Not a field delimiter: continuation of a multi-line value.
			if currentKey != "" {
				currentVal.WriteString("\n")
				currentVal.WriteString(line)
			}
			continue
		}

		// Line starts with "|": new field, flush the previous one first.
		flush()

		trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "|"))

		eqIdx := strings.Index(trimmed, "=")
		// No eqIdx means the wikitext is badly formatted, go to next line
		if eqIdx == -1 {
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		val := strings.TrimSpace(trimmed[eqIdx+1:])

		if key == "" {
			continue
		}

		currentKey = key
		currentVal.WriteString(val)
	}
	flush()

	return result
}

// ParseInfobox extracts key/value pairs from a wikitext infobox.
func ParseInfobox(wikitext string) map[string]string {
	fields := make(map[string]string)
	raw := extractRawPairs(wikitext)
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
