package main

import (
	"cmp"
	"errors"
	"slices"
	"strings"
	"unicode"
)

const (
	asc sortOrder = true
	desc sortOrder = false
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

func justTheKeys[E cmp.Ordered](index map[E]E) []E {
	var keys []E

	if len(index) == 0 {
		return keys
	}

	for k := range index {
		keys = append(keys, k)
	}
	return keys
}

func sortSlices[E cmp.Ordered](s []E, order sortOrder) []E {
	if len(s) == 0 {
		return s
	}

	if !order {
		slices.SortStableFunc(s, func(i, j E) int {
			return cmp.Compare(j, i)
		})
		return s
	}
	slices.SortStableFunc(s, func(i, j E) int {
		return cmp.Compare(i, j)
	})
	return s
}

func boolToPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func isNonZero(value any) (any, *bool, error) {
	switch v := value.(type) {
	case int:
		return value, new(v != 0), nil 
	case int32:
		return value, new(v != 0), nil
	case float32:
		return value, new(v != 0), nil
	case float64:
		return value, new(v != 0), nil
	case string:
		return value, new(v != ""), nil
	case bool:
		return value, new(v), nil
	case *int:
		return value, new(v != nil), nil
	case *int32:
		return value, new(v != nil), nil
	case *float32:
		return value, new(v != nil), nil
	case *float64:
		return value, new(v != nil), nil
	case *string:
		return value, new(v != nil), nil
	case *bool:
		return value, new(v != nil), nil
	case []string:
		return value, new(v != nil), nil
	case []int:
		return value, new(v != nil), nil
	case []bool:
		return value, new(v != nil), nil
	case map[string]string:
		return value, new(v != nil), nil
	case map[string]int:
		return value, new(v != nil), nil
	case map[int]int:
		return value, new(v != nil), nil
	case map[int]string:
		return value, new(v != nil), nil
	case map[string]bool:
		return value, new(v != nil), nil
	case map[int]bool:
		return value, new(v != nil), nil
	default:
		return value, nil, errors.New("unsupported type")
	}
}

func assignNonZero(a, b *string) {
	if a == nil || b == nil {
		return
	}
	_, ok, err := isNonZero(*a)
	if ok == nil || !*ok || err != nil {
		return
	}
	_, ok, err = isNonZero(*a)
	if ok == nil || !*ok || err != nil {
		return
	}
	*a = *b
}
