package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
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

func justTheKeys[E cmp.Ordered](index map[E]E) ([]E, error) {
	var keys []E

	if len(index) == 0 {
		return keys, errors.New("map is empty")
	}

	for k := range index {
		keys = append(keys, k)
	}
	return keys, nil
}

func sortSlices[E cmp.Ordered](s []E, asc bool) ([]E, error) {
	if len(s) == 0 {
		return s, errors.New("unable to sort empty slice")
	}

	if !asc {
		slices.SortStableFunc(s, func(i, j E) int {
			return cmp.Compare(j, i)
		})
		return s, nil
	}
	slices.SortStableFunc(s, func(i, j E) int {
		return cmp.Compare(i, j)
	})
	return s, nil
}

func boolToPtr(b bool) *bool {
	return &b
}

func boolToPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func isNonZero(value any) (*bool, string, error) {
	typeName := fmt.Sprintf("%T", value)
	switch v := value.(type) {
	case int:
		return boolToPtr(v != 0), typeName, nil
	case int32:
		return boolToPtr(v != 0), typeName, nil
	case float32:
		return boolToPtr(v != 0), typeName, nil
	case float64:
		return boolToPtr(v != 0), typeName, nil
	case string:
		return boolToPtr(v != ""), typeName, nil
	case bool:
		return boolToPtr(v), typeName, nil
	case *int:
		return boolToPtr(v != nil), typeName, nil
	case *int32:
		return boolToPtr(v != nil), typeName, nil
	case *float32:
		return boolToPtr(v != nil), typeName, nil
	case *float64:
		return boolToPtr(v != nil), typeName, nil
	case *string:
		return boolToPtr(v != nil), typeName, nil
	case *bool:
		return boolToPtr(v != nil), typeName, nil
	case []string:
		return boolToPtr(v != nil), typeName, nil
	case []int:
		return boolToPtr(v != nil), typeName, nil
	case []bool:
		return boolToPtr(v != nil), typeName, nil
	case map[string]string:
		return boolToPtr(v != nil), typeName, nil
	case map[string]int:
		return boolToPtr(v != nil), typeName, nil
	case map[int]int:
		return boolToPtr(v != nil), typeName, nil
	case map[int]string:
		return boolToPtr(v != nil), typeName, nil
	case map[string]bool:
		return boolToPtr(v != nil), typeName, nil
	case map[int]bool:
		return boolToPtr(v != nil), typeName, nil
	default:
		err := errors.New("unrecognized type")
		return nil, typeName, err
	}
}
