package ui

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

/*
	Holds each segment on the CLI, e.g. [ 10.12.90.12 ][ gobar ]
*/
type PromptSegment struct {
	Text    string
	Fgcolor string
	Bgcolor string
}

/*
	Internal representation for the current state of the command line.
	input: 		Holds user input (printable characters)
	escape: 	the current escape sequence
	cursor: 	Index of the current cursor
	tabCont: 	How many times we've hit Tab (mod 2)

*/

type commandLine struct {
	input    string
	escape   string
	cursor   int
	tabCount int
}

type history struct {
	commandHistory []string
	index          int
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
	EOT   = "\x04"

	PRINTABLE = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~`!@#$%^&*()_-+='\"{}[]\\|:;<>,. "

	historySize = 100
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

	prepared       = false
	commandHistory = newHistory(historySize)
	historyIndex   = 0

	override_colors = false
	override_fg     = ""
	override_bg     = ""
)

// Exported Functions

func DisplayPrompt(segments []PromptSegment) {
	for i, segment := range segments {
		fmt.Print(renderSegment(segment, i+1 == len(segments)))
	}
	color.Unset()
	fmt.Print(" ")
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
	resetKeyboard()
}

func Error(message string, uiSegments []PromptSegment) {
	override_colors = true
	override_fg = "black"
	override_bg = "red"
	DisplayPrompt(uiSegments)
	override_colors = false

	fmt.Println(message)
}

func Output(message string, uiSegments []PromptSegment) {
	override_colors = false
	DisplayPrompt(uiSegments)
	fmt.Println(message)
}

// Color functions

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

// Render and display functions

func renderSegment(segment PromptSegment, end bool) string {
	text := segment.Text
	fg := segment.Fgcolor
	bg := segment.Bgcolor
	if override_colors {
		fg = override_fg
		bg = override_bg
	}
	if end {
		text = colorize(text+" ", fg, bg)
		text += colorize(*prompt_end, bg, fg)
	} else {
		text = colorize(text+*prompt_mid, fg, bg)
	}
	return text
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

// Input processing functions

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

	if len(line.input) != 0 {
		commandHistory.push(line.input)
	}
	return line.input
}

// Returns -1 if no escape sequence found
func hasEscape(input string) int {
	return strings.Index(input, ESCSEQ)
}

// Returns the found escape sequence, empty string if none found
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

// Process escape sequences
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
	case HOME:
		line.cursor = 0
	case END:
		line.cursor = len(line.input)
	case UP:
		prev := commandHistory.previous()
		if prev != "" {
			line.input = prev
		}
		line.cursor = len(line.input)
	case DOWN:
		next := commandHistory.next()
		line.input = next
		line.cursor = len(line.input)
	}
}

// Process input after an ESC character is found
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

// Special characters such as tabs, backspaces etc
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
	case 0x0c: // Ctrl + l
		clearScreen()
	case 0x12: // Ctrl + r
		//clearScreen()
	case 0x04: // Ctrl + d
		line.input = ""
		line.cursor = 0
	}
}

// Returns true when line is complete
func (line *commandLine) handleInput(input byte, tabComplete func(string, int) string) bool {

	//debug("Key: %v", input)

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

// Utility functions

// Convert a byte to a string type
func byteToString(b byte) string {
	buf := make([]byte, 1)
	buf[0] = b
	return string(buf)
}

// Delete a character from a string at a given index
func deleteChar(input string, index int) string {
	return input[:index] + input[index+1:]
}

// Insert a character into a string at a given index
func insertChar(input string, index int, char byte) string {
	if len(input) == 0 {
		return byteToString(char)
	}
	return (input[:index] + byteToString(char) + input[index:])
}

// TTY functions

// Disables echo, disables newline buffering
func prepareKeyboard() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1", "-echo").Run()
	prepared = true
}

// Enables echo and sane settings
func resetKeyboard() {
	exec.Command("stty", "-F", "/dev/tty", "icanon", "sane").Run()
}

// Clear the screen, reposition to top of screen
func clearScreen() {
	// https://stackoverflow.com/questions/10105666/clearing-the-terminal-screen#15559322
	fmt.Print(ESCSEQ + "2J" + END)
	fmt.Print(ESCSEQ + "H" + END)
}

///// History Methods \\\\\

func (hist *history) isFull() bool {
	debug("Full? %d == %d", len(hist.commandHistory), cap(hist.commandHistory))
	return len(hist.commandHistory) == cap(hist.commandHistory)
}

// Push to history, returns index of last pushed item
func (hist *history) push(command string) int {
	// If we reused a history item
	if hist.index < len(hist.commandHistory) {
		hist.reuse(hist.index)
		hist.index = len(hist.commandHistory)
		return hist.index
	}

	if hist.isFull() {
		// Create new slice
		new := make([]string, 0, cap(hist.commandHistory))
		// Copy history into it TODO: more efficent way?
		new = append(new, hist.commandHistory[1:]...)

		hist.commandHistory = append(new, command)
	} else {
		hist.commandHistory = append(hist.commandHistory, command)
	}
	hist.index = len(hist.commandHistory)
	return len(hist.commandHistory) - 1
}

// Moves a history item at 'index' to the end
func (hist *history) reuse(index int) int {
	if index >= len(hist.commandHistory) {
		return -1
	}
	command := hist.commandHistory[index]
	new := make([]string, 0, cap(hist.commandHistory))
	new = append(new, hist.commandHistory[:index]...)
	new = append(new, hist.commandHistory[index+1:]...)
	new = append(new, command)
	hist.commandHistory = new
	return len(hist.commandHistory) - 1
}

/*
Return the next history item
*/
func (hist *history) next() (next string) {
	if hist.index+1 < len(hist.commandHistory) {
		hist.index += 1
		if hist.index >= len(hist.commandHistory) {
			hist.index = len(hist.commandHistory) // always point to one beyond
			next = ""
		} else {
			next = hist.commandHistory[hist.index]
		}
	}
	return next
}

/*
Return the previous history item
*/
func (hist *history) previous() (prev string) {

	hist.index -= 1
	if hist.index < 0 {
		hist.index = 0
	}

	if hist.index < len(hist.commandHistory) {
		prev = hist.commandHistory[hist.index]
	}
	return prev
}

func newHistory(capacity int) history {
	return history{
		make([]string, 0, capacity),
		0,
	}
}

// Returns -1 if no match found, otherwise the index of the first occurance
func searchHistory(partial string, start int) int {
	return -1
}
