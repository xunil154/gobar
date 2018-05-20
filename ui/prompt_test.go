package ui

import (
	"testing"
)

func TestInsertChar(t *testing.T) {
	cases := []struct {
		target   string
		char     byte
		index    int
		expected string
	}{
		{"", 'A', 0, "A"},
		{"asdf", '1', 0, "1asdf"},
		{"asdf", '1', 1, "a1sdf"},
		{"asdf", '1', 2, "as1df"},
		{"asdf", '1', 3, "asd1f"},
		{"asdf", '1', 4, "asdf1"},
	}

	for _, c := range cases {
		got := insertChar(c.target, c.index, c.char)
		if got != c.expected {
			t.Errorf("insertChar(%q, %q, %q) == %q, want %q",
				c.target, c.index, c.char, got, c.expected)
		}
	}
}

func TestDeleteChar(t *testing.T) {
	cases := []struct {
		target   string
		index    int
		expected string
	}{
		{"asdf", 0, "sdf"},
		{"asdf", 1, "adf"},
		{"asdf", 2, "asf"},
		{"asdf", 3, "asd"},
	}

	for _, c := range cases {
		got := deleteChar(c.target, c.index)
		if got != c.expected {
			t.Errorf("insertChar(%q, %q) == %q, want %q",
				c.target, c.index, got, c.expected)
		}
	}
}
