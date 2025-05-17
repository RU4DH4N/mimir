package helper

import (
	"strings"
	"unicode"
)

// need to unit test this
func Slugify(text string) string {
	var sb strings.Builder
	text = strings.ToLower(text)

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			sb.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' {
			if sb.Len() > 0 && sb.String()[sb.Len()-1] != '-' {
				sb.WriteRune('-')
			}
		} else {
			if sb.Len() > 0 && sb.String()[sb.Len()-1] != '-' {
				sb.WriteRune('-')
			}
		}
	}

	return strings.Trim(sb.String(), "-")
}

// this is temporary
func Linkify(text string) string {
	parts := strings.Split(text, "/")
	var sb strings.Builder
	for _, part := range parts {
		sb.WriteString("/" + Slugify(part))
	}
	return sb.String()
}
