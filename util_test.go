package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("actual length (%d) does not match expected length (%d)", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("actual word (%s) does not match expected word (%s)", word, expectedWord)
			}
		}
	}
}

func TestSanitizeInput(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "Ash Ketchum",
			expected: "ash_ketchum",
		},
		{
			input:    "Gotta Catch 'Em All",
			expected: "gotta_catch__em_all",
		},
		{
			input:    "M!sty",
			expected: "m_sty",
		},
		{
			input:    "Br()(k",
			expected: "br___k",
		},
	}

	for _, c := range cases {
		actual := sanitizeInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("actual length (%d) does not match expected length (%d)", len(actual), len(c.expected))
		}
		if actual != c.expected {
			t.Errorf("received %s, expected %s", actual, c.expected)
		}
	}
}

func TestCapitalize(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "lotad",
			expected: "Lotad",
		},
		{
			input:    "This is pokemon",
			expected: "This Is Pokemon",
		},
	}

	for _, c := range cases {
		actual := capitalize(c.input)
		if actual != c.expected {
			t.Errorf("received %s, expected %s", actual, c.expected)
		}
	}
}

func TestIsNonZero(t *testing.T) {
	type expected struct {
		nonZero  *bool
		typeName string
		err      error
	}

	cases := []struct {
		input    any
		expected expected
	}{
		{
			input: "",
			expected: expected{
				nonZero:  boolToPtr(false),
				typeName: "string",
				err: nil,
			},
		},
		{
			input: "test",
			expected: expected{
				nonZero:  boolToPtr(true),
				typeName: "string",
				err: nil,
			},
		},
		{
			input: 1,
			expected: expected{
				nonZero:  boolToPtr(true),
				typeName: "int",
				err: nil,
			},
		},
		{
			input: true,
			expected: expected{
				nonZero:  boolToPtr(true),
				typeName: "bool",
				err: nil,
			},
		},
		{
			input: []int{1, 2, 3},
			expected: expected{
				nonZero:  boolToPtr(true),
				typeName: "[]int",
				err: nil,
			},
		},
		{
			input: map[string]string{
				"Ash Ketchum": "ash_ketchum",
				"Gary Oak":    "gary_oak",
				"Red":         "red",
			},
			expected: expected{
				nonZero:  boolToPtr(true),
				typeName: "map[string]string",
				err: nil,
			},
		},
	}

	for _, c := range cases {
		actualNonZero, actualTypeName, actualErr := isNonZero(c.input)
		expected := c.expected

		if !boolToPtrEqual(actualNonZero, expected.nonZero) {
			t.Errorf("actual nonZero (%v) does not match expected nonZero (%v)", *actualNonZero, *expected.nonZero)
		} 
		if actualTypeName != expected.typeName {
			t.Errorf("actual typeName (%s) does not match expected typeName (%s)", actualTypeName, expected.typeName)
		}
		if actualErr != expected.err {
			t.Errorf("actual err (%v) does not match expected err (%v)", actualErr, expected.err)
		}
	}
}
