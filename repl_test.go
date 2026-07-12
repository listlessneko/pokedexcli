package main

import (
	"testing"
	"bytes"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "  hello  world  ",
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

func TestCommandInspect(t *testing.T) {
	cfg := &config {
		Caught: map[string]Pokemon {
			"lotad": Pokemon {
				Name: "lotad",
				Height: 5,
				Weight: 26,
				Stats: []PokemonStats {
					{
						BaseStat: 40,
						Stat: NamedAPIResource {
							Name: "hp",
						},
					},
				},
				Types: []PokemonTypes {
					{
						Type: NamedAPIResource {
							Name: "water",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := commandInspect(cfg, &buf, []string{"lotad"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Name: lotad\nHeight: 5\nWeight: 26\nStats:\n- hp: 40\nTypes:\n- water\n"

	if buf.String() != expected {
		t.Errorf("received %q, expected %q", buf.String(), expected)
	}
}

