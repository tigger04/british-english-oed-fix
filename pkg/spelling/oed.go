// ABOUTME: OED spelling engine. Loads embedded word lists and performs
// case-preserving whole-word replacement on input text.
package spelling

import (
	"strings"
	"unicode"
)

// OEDEngine holds the merged lookup map and tracks replacement counts.
type OEDEngine struct {
	words   map[string]string
	Changes int
}

// NewOEDEngine creates an engine from the provided word list data strings.
func NewOEDEngine(wordLists ...string) (*OEDEngine, error) {
	e := &OEDEngine{
		words: make(map[string]string),
	}
	for _, data := range wordLists {
		if err := e.loadWordList(data); err != nil {
			return nil, err
		}
	}
	return e, nil
}

// loadWordList parses lines of "wrong=correct" into the lookup map.
func (e *OEDEngine) loadWordList(data string) error {
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key != "" && val != "" {
			e.words[strings.ToLower(key)] = strings.ToLower(val)
		}
	}
	return nil
}

// ProcessLine replaces words in a single line, preserving case.
func (e *OEDEngine) ProcessLine(line string) string {
	runes := []rune(line)
	var result strings.Builder
	result.Grow(len(line))

	i := 0
	for i < len(runes) {
		if isWordChar(runes[i]) {
			// Extract the whole word
			j := i
			for j < len(runes) && isWordChar(runes[j]) {
				j++
			}
			word := string(runes[i:j])
			replaced := e.replaceWord(word)
			result.WriteString(replaced)
			i = j
		} else {
			result.WriteRune(runes[i])
			i++
		}
	}
	return result.String()
}

// isWordChar returns true for letters and apostrophes (word-internal).
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || r == '\''
}

// replaceWord looks up a word and applies case-preserving replacement.
func (e *OEDEngine) replaceWord(word string) string {
	lower := strings.ToLower(word)
	replacement, ok := e.words[lower]
	if !ok {
		return word
	}
	e.Changes++
	return applyCase(word, replacement)
}

// applyCase transfers the case pattern of orig onto replacement.
func applyCase(orig, replacement string) string {
	if orig == strings.ToLower(orig) {
		return replacement
	}
	if orig == strings.ToUpper(orig) {
		return strings.ToUpper(replacement)
	}
	// Title Case: first letter uppercase, rest lowercase
	origRunes := []rune(orig)
	if len(origRunes) > 0 && unicode.IsUpper(origRunes[0]) {
		runes := []rune(replacement)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}
	// Mixed case fallback: return lowercase
	return replacement
}
