package helper

import (
	"regexp"
	"strings"
)

// Regex (See: https://regexlicensing.org/)

// Remove all non-word chars (fix for UTF-8 chars)
// allows unicode letters & numbers, and also '-'
// \p{L} matches any kind of letter from any language
// \p{N} matches any kind of numeric character in any script
var reNonWord = regexp.MustCompile(`[^\p{L}\p{N}\-]`)

// Replace multiple - with single -
var reMultiDash = regexp.MustCompile(`-{2,}`)

func Slugify(text string) string {
	text = strings.ToLower(text)

	text = reNonWord.ReplaceAllString(text, "-")

	text = reMultiDash.ReplaceAllString(text, "-")

	text = strings.TrimLeft(text, "-")
	text = strings.TrimRight(text, "-")

	return text
}
