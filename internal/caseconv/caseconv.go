// Package caseconv converts strings between kebab, snake, camel, pascal, and title case.
//
// Word boundaries: input is trimmed, split on spaces/underscores/hyphens, and on
// CamelCase boundaries (e.g. "HelloWorld" → "Hello", "World"). Multiple spaces
// and leading/trailing spaces are normalized. Accented characters are normalized
// to ASCII equivalents (e.g. é→e, í→i). Case conversion is Unicode-aware.
package caseconv

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Words splits s into words: trims space, optionally normalizes accented chars
// to ASCII (when raw is false), splits on spaces/underscores/hyphens, and on
// CamelCase boundaries.
func Words(s string, raw bool) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !raw {
		s = normalizeAccents(s)
	}

	// Replace delimiters with space so we have one splitting rule
	s = replaceDelimiters(s)

	var words []string
	var buf strings.Builder

	runes := []rune(s)
	for i, r := range runes {
		if r == ' ' || r == '\t' {
			if buf.Len() > 0 {
				words = append(words, buf.String())
				buf.Reset()
			}
			continue
		}

		// Case boundary: new word before uppercase when
		// - previous was lowercase, or
		// - previous was uppercase and current is uppercase followed by lowercase (e.g. "HTTP" before "Response")
		if buf.Len() > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsLower(prev) || (unicode.IsUpper(prev) && nextLower) {
				words = append(words, buf.String())
				buf.Reset()
			}
		}

		buf.WriteRune(r)
	}

	if buf.Len() > 0 {
		words = append(words, buf.String())
	}

	return words
}

// normalizeAccents converts accented Unicode to ASCII (é→e, í→i, ñ→n, etc.)
// using NFD decomposition and stripping combining marks.
func normalizeAccents(s string) string {
	nfd := norm.NFD.String(s)
	var b strings.Builder
	for _, r := range nfd {
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) || unicode.Is(unicode.Me, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func replaceDelimiters(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '_' || r == '-' {
			b.WriteRune(' ')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ToKebab returns lowercase, hyphen-separated (hello-world). When raw is true, accented chars are preserved.
func ToKebab(s string, raw bool) string {
	words := Words(s, raw)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "-")
}

// ToSnake returns lowercase, underscore-separated (hello_world). When raw is true, accented chars are preserved.
func ToSnake(s string, raw bool) string {
	words := Words(s, raw)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

// ToCamel returns camelCase (helloWorld): first word lower, rest title-cased. When raw is true, accented chars are preserved.
func ToCamel(s string, raw bool) string {
	words := Words(s, raw)
	if len(words) == 0 {
		return ""
	}
	for i, w := range words {
		if i == 0 {
			words[i] = strings.ToLower(w)
		} else {
			words[i] = toTitleWord(w)
		}
	}
	return strings.Join(words, "")
}

// ToPascal returns PascalCase (HelloWorld): each word title-cased. When raw is true, accented chars are preserved.
func ToPascal(s string, raw bool) string {
	words := Words(s, raw)
	for i, w := range words {
		words[i] = toTitleWord(w)
	}
	return strings.Join(words, "")
}

// ToTitle returns Title Case (Hello World): each word title-cased, space-separated. When raw is true, accented chars are preserved.
func ToTitle(s string, raw bool) string {
	words := Words(s, raw)
	for i, w := range words {
		words[i] = toTitleWord(w)
	}
	return strings.Join(words, " ")
}

// toTitleWord uppercases the first rune, lowercases the rest (e.g. "hello" → "Hello").
func toTitleWord(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// Detect returns the case style of s: "kebab", "snake", "camel", "pascal", "title", or "unknown".
func Detect(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown"
	}

	hasSpace := strings.Contains(s, " ")
	hasUnderscore := strings.Contains(s, "_")
	hasHyphen := strings.Contains(s, "-")

	if hasSpace {
		return "title"
	}
	if hasUnderscore && hasHyphen {
		return "unknown"
	}
	if hasUnderscore {
		return "snake"
	}
	if hasHyphen {
		return "kebab"
	}

	// No separators: camel vs pascal by first rune
	runes := []rune(s)
	if len(runes) == 0 {
		return "unknown"
	}
	if unicode.IsLower(runes[0]) {
		return "camel"
	}
	return "pascal"
}
