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

func TestHistoryPush(t *testing.T) {
	hist := newHistory(5)
	cases := []struct {
		command        string
		expected_index int
		first_command  string
	}{
		{"A", 0, "A"},
		{"B", 1, "A"},
		{"C", 2, "A"},
		{"D", 3, "A"},
		{"E", 4, "A"},
		{"F", 4, "B"},
		{"G", 4, "C"},
	}

	for _, c := range cases {
		got := hist.push(c.command)
		if got != c.expected_index {
			t.Errorf("history.push(%q) == %q, want %q",
				c.command, got, c.expected_index)
		}

		first := hist.commandHistory[0]

		if first != c.first_command {
			t.Errorf("history.commandHistory[0] == %q, want %q",
				first, c.first_command)
		}
		if cap(hist.commandHistory) != 5 {
			t.Errorf("history.commandHistory did not maintain max capacity")
		}
	}
}

func TestLineHandleInput(t *testing.T) {
	line := newLine()
	tc := func(partial string, count int) string { return "" }

	//finished := line.handleInput(buf[0], tabComplete)
	cases := []struct {
		char     byte
		expected string
		finished bool
	}{
		{'A', "A", false},
		{'B', "AB", false},
		{'C', "ABC", false},
		{'!', "ABC!", false},
		{' ', "ABC! ", false},
		{0x7f, "ABC!", false}, // backspace
		{0x04, "", false},     // Ctrl + D
		{'A', "A", false},
		{'\n', "A", true},
	}

	for _, c := range cases {
		got := line.handleInput(c.char, tc)

		if got != c.finished {
			t.Errorf("commandLine.handleInput(%q, tc) == %v, want %v",
				c.char, got, c.finished)
		}
		if line.input != c.expected {
			t.Errorf("commandLine.handleInput.input = %q, want %q",
				line.input, c.expected)
		}
	}
}

func TestHistoryReuse(t *testing.T) {
	hist := newHistory(5)

	commands := []string{"A", "B", "C", "D"}

	for _, c := range commands {
		hist.push(c)
	}

	cases := []struct {
		index          int
		expected_index int
		first_command  string
		last_command   string
	}{
		{0, 3, "B", "A"}, // BCDA
		{0, 3, "C", "B"}, // CDAB
		{1, 3, "C", "D"}, // CABD
	}

	for _, c := range cases {
		got := hist.reuse(c.index)

		if got != c.expected_index {
			t.Errorf("history.reuse(%q) == %q, want %q",
				c.index, got, c.expected_index)
		}

		first := hist.commandHistory[0]
		last := hist.commandHistory[len(hist.commandHistory)-1]

		if first != c.first_command {
			t.Errorf("history.commandHistory[0] == %q, want %q",
				first, c.first_command)
		}
		if last != c.last_command {
			t.Errorf("history.commandHistory[-1] == %q, want %q",
				last, c.last_command)
		}
		if cap(hist.commandHistory) != 5 {
			t.Errorf("history.commandHistory did not maintain max capacity")
		}
	}
}
