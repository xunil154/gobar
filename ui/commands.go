package ui

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type CommandOutput struct {
	Command   string
	StartTime time.Time
	EndTime   time.Time
	Time      time.Duration
	Output    string
	Error     bool
}

var (
	commands = make(map[string](func(string, *[]PromptSegment) (string, error)))
)

func RegisterCommand(name string, callback func(string, *[]PromptSegment) (string, error)) {
	commands[name] = callback
}

func isValidCommand(command string) bool {
	for name, _ := range commands {
		if name == command {
			return true
		}
	}
	return false
}

func TabComplete(partial string, tabcount int) string {
	return partial
}

func ProcessInput(command string, uiSegments *[]PromptSegment) (output CommandOutput, err error) {
	initCommands()

	args := strings.Fields(command)
	if len(args) == 0 {
		return output, nil // do nothing
	}

	output.StartTime = time.Now()
	if isValidCommand(args[0]) {
		subcmd := strings.Join(args[1:], " ")
		// Call function with arguments
		output.Output, err = commands[args[0]](subcmd, uiSegments)
	} else {
		output.Output, err = help(command, uiSegments)
	}

	output.EndTime = time.Now()
	output.Time = output.EndTime.Sub(output.StartTime)

	return output, err
}

// Commands
func initCommands() {
	if len(commands) == 0 {
		RegisterCommand("help", help)
		RegisterCommand("exit", exit)
	}
}

func exit(command string, uiSegments *[]PromptSegment) (string, error) {
	// Do nothing really, handeled in gobar.go
	// TODO: Make this cleaner
	return "", errors.New("exit")
}

func help(command string, uiSegments *[]PromptSegment) (string, error) {
	args := strings.Fields(command)
	if len(args) > 0 && !isValidCommand(args[0]) {
		args := strings.Split(command, " ")
		return "", errors.New(fmt.Sprintf("Command '%v' not found", args[0]))
	}
	var help string = "Available commands:\n"
	for cmd, _ := range commands {
		help += fmt.Sprintf("\t%v\n", cmd)
	}
	return help, nil
}
