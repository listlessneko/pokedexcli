package main

import (
	"cmp"
	"errors"
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

func justTheKeys[K cmp.Ordered, V any](index map[K]V) ([]K, error) {
	var keys []K

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

func isNonZero(value any) (any, *bool, error) {
	switch v := value.(type) {
	case int:
		return value, boolToPtr(v != 0), nil
	case int32:
		return value, boolToPtr(v != 0), nil
	case float32:
		return value, boolToPtr(v != 0), nil
	case float64:
		return value, boolToPtr(v != 0), nil
	case string:
		return value, boolToPtr(v != ""), nil
	case bool:
		return value, boolToPtr(v), nil
	case *int:
		return value, boolToPtr(v != nil), nil
	case *int32:
		return value, boolToPtr(v != nil), nil
	case *float32:
		return value, boolToPtr(v != nil), nil
	case *float64:
		return value, boolToPtr(v != nil), nil
	case *string:
		return value, boolToPtr(v != nil), nil
	case *bool:
		return value, boolToPtr(v != nil), nil
	case []string:
		return value, boolToPtr(v != nil), nil
	case []int:
		return value, boolToPtr(v != nil), nil
	case []bool:
		return value, boolToPtr(v != nil), nil
	case map[string]string:
		return value, boolToPtr(v != nil), nil
	case map[string]int:
		return value, boolToPtr(v != nil), nil
	case map[int]int:
		return value, boolToPtr(v != nil), nil
	case map[int]string:
		return value, boolToPtr(v != nil), nil
	case map[string]bool:
		return value, boolToPtr(v != nil), nil
	case map[int]bool:
		return value, boolToPtr(v != nil), nil
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
	return
}

func newTrie() *trie {
	return &trie{Root: &trieNode{}}
}

func (t *trie) add(word string) {
	currentLevel := t.Root
	for _, letter := range word {
		if currentLevel.Children == nil {
			currentLevel.Children = make(children)
		}
		if currentLevel.Children[letter] == nil {
			currentLevel.Children[letter] = &trieNode{}
		}
		currentLevel = currentLevel.Children[letter]
	}
	currentLevel.End = true
}

func (t *trie) searchLevel(currentLevel *trieNode, currentPrefix string, words []string) []string {
	if currentLevel.End {
		words = append(words, currentPrefix)
	}
	keys, err := justTheKeys(currentLevel.Children)
	if err == nil {
		return words
	}
	sortedKeys, err := sortSlices(keys, true)
	if err == nil {
		return words
	}
	for _, letter := range sortedKeys {
		if !currentLevel.End {
			t.searchLevel(currentLevel.Children[letter], currentPrefix + string(letter), words)
		}
	}
	return words
}

func (t *trie) searchByPrefix(prefix string) []string {
	var collectedWords []string
	currentLevel := t.Root
	for _, letter := range prefix {
		if currentLevel.Children == nil || currentLevel.Children[letter] == nil {
			return collectedWords
		}
		currentLevel = currentLevel.Children[letter]
	}
	return t.searchLevel(currentLevel, prefix, collectedWords)
}
