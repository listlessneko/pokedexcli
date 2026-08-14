package main

// TODO: Should rewrite due to refactor.

// import (
// 	"reflect"
// 	"testing"
// )
//
// func TestCleanInput(t *testing.T) {
// 	cases := []struct {
// 		input    string
// 		expected []string
// 	}{
// 		{
// 			input:    "  hello  world  ",
// 			expected: []string{"hello", "world"},
// 		},
// 	}
//
// 	for _, c := range cases {
// 		actual := cleanInput(c.input)
//
// 		if len(actual) != len(c.expected) {
// 			t.Errorf("actual length (%d) does not match expected length (%d)", len(actual), len(c.expected))
// 		}
//
// 		for i := range actual {
// 			word := actual[i]
// 			expectedWord := c.expected[i]
//
// 			if word != expectedWord {
// 				t.Errorf("actual word (%s) does not match expected word (%s)", word, expectedWord)
// 			}
// 		}
// 	}
// }
//
// func TestSanitizeInput(t *testing.T) {
// 	cases := []struct {
// 		input    string
// 		expected string
// 	}{
// 		{
// 			input:    "",
// 			expected: "",
// 		},
// 		{
// 			input:    "Ash Ketchum",
// 			expected: "ash_ketchum",
// 		},
// 		{
// 			input:    "Gotta Catch 'Em All",
// 			expected: "gotta_catch__em_all",
// 		},
// 		{
// 			input:    "M!sty",
// 			expected: "m_sty",
// 		},
// 		{
// 			input:    "Br()(k",
// 			expected: "br___k",
// 		},
// 	}
//
// 	for _, c := range cases {
// 		actual := sanitizeInput(c.input)
//
// 		if len(actual) != len(c.expected) {
// 			t.Errorf("actual length (%d) does not match expected length (%d)", len(actual), len(c.expected))
// 		}
// 		if actual != c.expected {
// 			t.Errorf("received %s, expected %s", actual, c.expected)
// 		}
// 	}
// }
//
// func TestCapitalize(t *testing.T) {
// 	cases := []struct {
// 		input    string
// 		expected string
// 	}{
// 		{
// 			input:    "",
// 			expected: "",
// 		},
// 		{
// 			input:    "lotad",
// 			expected: "Lotad",
// 		},
// 		{
// 			input:    "This is pokemon",
// 			expected: "This Is Pokemon",
// 		},
// 	}
//
// 	for _, c := range cases {
// 		actual := capitalize(c.input)
// 		if actual != c.expected {
// 			t.Errorf("received %s, expected %s", actual, c.expected)
// 		}
// 	}
// }
//
// func TestIsNonZero(t *testing.T) {
// 	type expected struct {
// 		value   any
// 		boolPtr *bool
// 		err     error
// 	}
//
// 	cases := []struct {
// 		input    any
// 		expected expected
// 	}{
// 		{
// 			input: "",
// 			expected: expected{
// 				value:   "",
// 				boolPtr: boolToPtr(false),
// 				err:     nil,
// 			},
// 		},
// 		{
// 			input: "test",
// 			expected: expected{
// 				value:   "test",
// 				boolPtr: boolToPtr(true),
// 				err:     nil,
// 			},
// 		},
// 		{
// 			input: 1,
// 			expected: expected{
// 				value:   1,
// 				boolPtr: boolToPtr(true),
// 				err:     nil,
// 			},
// 		},
// 		{
// 			input: true,
// 			expected: expected{
// 				value:   true,
// 				boolPtr: boolToPtr(true),
// 				err:     nil,
// 			},
// 		},
// 		{
// 			input: []int{1, 2, 3},
// 			expected: expected{
// 				value:   []int{1, 2, 3},
// 				boolPtr: boolToPtr(true),
// 				err:     nil,
// 			},
// 		},
// 		{
// 			input: map[string]string{
// 				"Ash Ketchum": "ash_ketchum",
// 				"Gary Oak":    "gary_oak",
// 				"Red":         "red",
// 			},
// 			expected: expected{
// 				value: map[string]string{
// 					"Ash Ketchum": "ash_ketchum",
// 					"Gary Oak":    "gary_oak",
// 					"Red":         "red",
// 				},
// 				boolPtr: boolToPtr(true),
// 				err:      nil,
// 			},
// 		},
// 	}
//
// 	for _, c := range cases {
// 		actualValue, actualBoolPtr, actualErr := isNonZero(c.input)
// 		expected := c.expected
//
// 		if !reflect.DeepEqual(actualValue, expected.value) {
// 			t.Errorf("actual value (%v) does not match expected value (%v)", *actualBoolPtr, *expected.boolPtr)
// 		}
// 		if !boolToPtrEqual(actualBoolPtr, expected.boolPtr) {
// 			t.Errorf("actual bool (%v) does not match expected bool (%v)", *actualBoolPtr, *expected.boolPtr)
// 		}
// 		if actualErr != expected.err {
// 			t.Errorf("actual err (%v) does not match expected err (%v)", actualErr, expected.err)
// 		}
// 	}
// }
