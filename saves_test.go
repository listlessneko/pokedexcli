package main

// TODO: Should rewrite due to refactor.

// import (
// 	"maps"
// 	"testing"
// )
//
// func TestSaveLoadIndex(t *testing.T) {
// 	cases := []struct {
// 		index    map[string]string
// 		expected map[string]string
// 	}{
// 		{
// 			index: map[string]string{
// 				"test":        "test",
// 				"Ask Ketchum": "ash_ketchum",
// 			},
// 			expected: map[string]string{
// 				"test":        "test",
// 				"Ask Ketchum": "ash_ketchum",
// 			},
// 		},
// 	}
//
// 	for _, c := range cases {
// 		err := saveIndex(c.index)
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 		actual, err := loadIndex()
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
//
// 		if !maps.Equal(actual, c.expected) {
// 			t.Errorf("received %v, expected %v", actual, c.expected)
// 		}
// 	}
// }
