package webserver

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func normalizar(texto string) string {
	var t string
	n := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ := transform.String(n, texto)

	for _, v := range s {
		if v == '.' {
			t += string(v)
			continue
		}
		if unicode.IsLetter(v) || unicode.IsDigit(v) {
			t += string(v)
		} else {
			t += "_"
		}
	}
	if t == "" {
		return texto
	}
	return strings.ToLower(t)
}
