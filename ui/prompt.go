package ui

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

type PromptSegment struct {
	Text    string
	Fgcolor string
	Bgcolor string
}

type commandLine struct {
	input    string
	escape   string
	cursor   int
	tabCount int
}

const (
	ESC    = "\x1b"
	ESCSEQ = ESC + "\x5b"
	HOME   = ESCSEQ + "\x31\x7e"
	INSERT = ESCSEQ + "\x32\x7e"
	DELETE = ESCSEQ + "\x33\x7e"
	END    = ESCSEQ + "\x34\x7e"

	UP    = ESCSEQ + "\x41"
	DOWN  = ESCSEQ + "\x42"
	RIGHT = ESCSEQ + "\x43"
	LEFT  = ESCSEQ + "\x44"

	PRINTABLE = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~`!@#$%^&*()_-+='\"{}[]\\|:;<>,. "
)

var (
	ALLESC = [...]string{HOME, INSERT, DELETE, END, UP, DOWN, LEFT, RIGHT}

	// Possible prompt options ⌲  ▶  ⌦ ⫸     
	prompt_end = flag.String("prompt_end", "", "Prompt character")
	prompt_mid = flag.String("prompt_mid", " ", "Prompt character")

	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	black   = color.New(color.FgBlack).SprintFunc()
	white   = color.New(color.FgWhite).SprintFunc()
	bgred   = color.New(color.BgRed).SprintFunc()
	bggreen = color.New(color.BgGreen).SprintFunc()
	bgblue  = color.New(color.BgBlue).SprintFunc()
	bgblack = color.New(color.BgBlack).SprintFunc()
	bgwhite = color.New(color.BgWhite).SprintFunc()

	prepared = false
)

func (segment *PromptSegment) String() string {
	return segment.Text + " fg: " + segment.Fgcolor + " bg: " + segment.Bgcolor
}

func colorize(text string, fg string, bg string) string {
	var colorized string = text
	switch fg {
	case "white":
		colorized = white(text)
	case "red":
		colorized = red(text)
	case "green":
		colorized = green(text)
	case "black":
		colorized = black(text)
	case "blue":
		colorized = blue(text)
	}

	switch bg {
	case "white":
		colorized = bgwhite(colorized)
	case "red":
		colorized = bgred(colorized)
	case "green":
		colorized = bggreen(colorized)
	case "black":
		colorized = bgblack(colorized)
	case "blue":
		colorized = bgblue(colorized)
	}
	return colorized
}

func renderSegment(segment PromptSegment, end bool) string {
	text := segment.Text
	if end {
		text = colorize(text+" ", segment.Fgcolor, segment.Bgcolor)
		text += colorize(*prompt_end, segment.Bgcolor, segment.Fgcolor)
	} else {
		text = colorize(text+*prompt_mid, segment.Fgcolor, segment.Bgcolor)
	}
	return text
}

func DisplayPrompt(segments []PromptSegment) {
	for i, segment := range segments {
		fmt.Print(renderSegment(segment, i+1 == len(segments)))
	}
	color.Unset()
	fmt.Print(" ")
}

func drawLine(line commandLine, prompts []PromptSegment) {
	DisplayPrompt(prompts)
	fmt.Print(line.input)
	if line.cursor < len(line.input) {
		fmt.Printf(ESCSEQ+"%vD", len(line.input)-line.cursor)
	}
}

func redrawLine(line commandLine, prompts []PromptSegment) {
	fmt.Print("\r")
	length := 0
	for i, segment := range prompts {
		length += len(renderSegment(segment, i+1 == len(prompts)))
	}
	for i := 0; i < length+len(line.input); i++ { // clear line
		fmt.Print(" ")
	}
	fmt.Print("\r")
	drawLine(line, prompts)
}

func prepareKeyboard() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1", "-echo").Run()
	prepared = true
}

// Returns -1 if no escape sequence found
func hasEscape(input string) int {
	return strings.Index(input, ESCSEQ)
}

func getEscape(input string) string {
	index := hasEscape(input)
	if index != -1 {
		for _, esc := range ALLESC {
			if len(input) >= index+len(esc) {
				test := input[index : index+len(esc)]
				if test == esc {
					return esc
				}
			}
		}
	}
	return ""
}

func (line *commandLine) handleEscape(escape string) {
	switch escape {
	case LEFT:
		if line.cursor > 0 {
			line.cursor -= 1
		}
	case RIGHT:
		if line.cursor < len(line.input) {
			line.cursor += 1
		}
	case DELETE:
		if line.cursor < len(line.input) {
			line.input = deleteChar(line.input, line.cursor)
		}
	}
}

func deleteChar(input string, index int) string {
	return input[:index] + input[index+1:]
}
func insertChar(input string, index int, char byte) string {
	if len(input) == 0 {
		return byteToString(char)
	}
	return (input[:index] + byteToString(char) + input[index:])
}

func byteToString(b byte) string {
	buf := make([]byte, 1)
	buf[0] = b
	return string(buf)
}

func (line *commandLine) handleEscapeInput(input byte) {

	char := byteToString(input)

	// Push the next escape character
	line.escape += char

	// Check if we match known escape sequences
	esc := getEscape(line.escape)
	if esc != "" {
		line.handleEscape(esc)
		line.escape = ""
	}
}

func (line *commandLine) handleSpecialInput(input byte, tabComplete func(string, int) string) {

	switch input {
	case 0x7F: // backspace
		if line.cursor > 0 {
			line.input = deleteChar(line.input, line.cursor-1)
			line.cursor -= 1 // Move cursor back one
		}
	case 0x09: // \t
		line.input = tabComplete(line.input, line.tabCount)
		line.tabCount = (line.tabCount + 1) % 2
	}
}

// Returns true when line is complete
func (line *commandLine) handleInput(input byte, tabComplete func(string, int) string) bool {

	char := byteToString(input)

	switch char {
	case "\n":
		return true // Done
	case ESC:
		line.escape = ESC
	}

	// If we are processing an escape sequence
	if len(line.escape) != 0 {
		line.handleEscapeInput(input)
		return false
	}

	// If we are processing normal printable characters
	if strings.Index(PRINTABLE, char) != -1 {
		line.input = insertChar(line.input, line.cursor, input)
		line.cursor += 1
	} else {
		line.handleSpecialInput(input, tabComplete)
	}
	return false
}

func getInput(prompts []PromptSegment, tabComplete func(string, int) string) string {
	line := commandLine{"", "", 0, 0}

	for {
		// read a single byte
		buf := make([]byte, 1)
		os.Stdin.Read(buf)

		// process the single byte
		finished := line.handleInput(buf[0], tabComplete)
		redrawLine(line, prompts)
		if finished {
			break
		}
	}

	return line.input
}

func GetUserInput(segments []PromptSegment, tabComplete func(string, int) string) string {
	if !prepared {
		prepareKeyboard()
	}
	DisplayPrompt(segments)
	//reader := bufio.NewReader(os.Stdin)
	//text, _ := reader.ReadString('\n')
	text := getInput(segments, tabComplete)
	text = strings.TrimSpace(text)
	return text
}

func Exit() {
	fmt.Println("") // newline to not mess up terminal
	info("Fixing terminal")
	exec.Command("stty", "-F", "/dev/tty", "icanon", "sane").Run()
}
