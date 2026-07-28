package main

import (
	"strings"
	"unicode"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func sanitizeInput(text string) string {
	if text == "" {
		return text
	}

	r := strings.NewReplacer(
		"`", "_",
		"~", "_",
		"!", "_",
		"#", "_",
		"$", "_",
		"%", "_",
		"^", "_",
		"&", "_",
		"*", "_",
		"(", "_",
		")", "_",
		"-", "_",
		"+", "_",
		"[", "_",
		"]", "_",
		"{", "_",
		"}", "_",
		"\\", "_",
		"|", "_",
		";", "_",
		":", "_",
		"'", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		",", "_",
		".", "_",
		"/", "_",
		"?", "_",
		" ", "_",
	)
	text = r.Replace(text)
	return strings.ToLower(text)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	final := ""
	for w := range strings.SplitSeq(strings.TrimSpace(s), " ") {
		r := []rune(w)
		r[0] = unicode.ToUpper(r[0])
		final = final + string(r) + " "
	}
	return strings.TrimSpace(final)
}
