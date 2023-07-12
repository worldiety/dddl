package model

import (
	"fmt"
	"strings"
	"unicode"
)

var symbolReplacementTable = map[rune]string{
	'Ü': "Ue",
	'Ä': "Ae",
	'Ö': "Oe",
	'ß': "ss",
	'@': "At",
	'ü': "ue",
	'ä': "ae",
	'ö': "oe",
}

func makePkgIdentifier(s string) string {
	return strings.ToLower(makeIdentifier(s))
}

func makeUpIdentifier(s string) string {
	ident := makeIdentifier(s)
	var sb strings.Builder
	for i, r := range ident {
		if i == 0 {
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// makeIdentifier uses the common rules of identifiers to create an acceptable string from an arbitrary input.
// This is suitable for c-like languages (Go, Java etc.)
func makeIdentifier(s string) string {
	if len(s) == 0 {
		return "_"
	}

	if isDigit(runeAt(s, 0)) {
		s = "_" + s
	}

	var sb strings.Builder
	var nextUp bool
	for _, r := range s {
		if rep, ok := symbolReplacementTable[r]; ok {
			sb.WriteString(rep)
		} else {
			if isLetter(r) {
				if nextUp {
					nextUp = false
					sb.WriteRune(unicode.ToUpper(r))
				} else {
					sb.WriteRune(r)
				}
			} else {
				nextUp = true
			}
		}
	}

	return sb.String()
}

func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func runeAt(s string, idx int) rune {
	for i, r := range s {
		if i == idx {
			return r
		}
	}

	panic(fmt.Sprintf("index out of bounds: max %d, idx %d", len(s), idx))
}
