package main

import (
	"bytes"
	"sort"
	"strings"
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

func TestCapitalize(t *testing.T) {
	cases := []struct {
		input string
		expected string
	}{
		{
			input: "",
			expected: "",
		},
		{
			input: "lotad",
			expected: "Lotad",
		},
	}

	for _, c := range cases {
		actual := capitalize(c.input)
		if actual != c.expected {
			t.Errorf("received %s, expected %s", actual, c.expected)
		}
	}
}

func TestCommandInspect(t *testing.T) {
	cases := []struct {
		args     []string
		caught   map[string]Pokemon
		expected string
	}{
		{
			args: []string{"lotad"},
			caught: map[string]Pokemon{
				"lotad": Pokemon{
					Name:   "lotad",
					Height: 5,
					Weight: 26,
					Stats: []PokemonStats{
						{
							BaseStat: 40,
							Stat: NamedAPIResource{
								Name: "hp",
							},
						},
					},
					Types: []PokemonTypes{
						{
							Type: NamedAPIResource{
								Name: "water",
							},
						},
					},
				},
			},
			expected: "Name: Lotad\nHeight: 5\nWeight: 26\nStats:\n- hp: 40\nTypes:\n- Water\n",
		},
		{
			args:     []string{"lotad"},
			caught:   map[string]Pokemon{},
			expected: "You have not caught that Pokemon.\n",
		},
		{
			args:     []string{},
			caught:   map[string]Pokemon{},
			expected: "Please provide a valid Pokemon.\n",
		},
	}

	for _, c := range cases {
		cfg := &config{Caught: c.caught}

		var buf bytes.Buffer
		err := commandInspect(cfg, &buf, c.args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if buf.String() != c.expected {
			t.Errorf("received %q, expected %q", buf.String(), c.expected)
		}
	}

}

func TestCommandPokedex(t *testing.T) {
	cases := []struct {
		caught   map[string]Pokemon
		expected []string
	}{
		{
			caught:   map[string]Pokemon{},
			expected: []string{"Your Pokedex:\nYou have no Pokemon...\n"},
		},
		{
			caught: map[string]Pokemon{
				"lotad": Pokemon{
					Name: "lotad",
				},
				"cyndaquil": Pokemon{
					Name: "cyndaquil",
				},
				"lugia": Pokemon{
					Name: "lugia",
				},
			},
			expected: []string{"- Cyndaquil", "- Lotad", "- Lugia"},
		},
	}

	for _, c := range cases {
		cfg := &config{Caught: c.caught}

		var buf bytes.Buffer
		err := commandPokedex(cfg, &buf, []string{""})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(c.caught) == 0 {
			if buf.String() != c.expected[0] {
				t.Errorf("received %v, expected %v", buf.String(), c.expected[0])
			}
			continue
		}

		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		caughtPokemon := lines[1:]
		sort.Strings(caughtPokemon)

		if len(caughtPokemon) != len(c.expected) {
			t.Errorf("actual length (%d) does not match expected length (%d)", len(caughtPokemon), len(c.expected))
			continue
		}

		for i, p := range caughtPokemon {
			if p != c.expected[i] {
				t.Errorf("received %v, expected %v", p, c.expected[i])
			}
		}
	}
}
