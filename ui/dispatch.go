package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type command struct {
	name        string
	description string
	help        string
	callback    func(input string) (string, error)
	tabComplete func(input string, tabcount int) string
}

func (cmd command) String() string {
	return fmt.Sprintf("Command %v: %v - %v (%v)(%v)",
		cmd.name, cmd.description, cmd.help, cmd.callback, cmd.tabComplete)
}

type CommandOutput struct {
	Command   string
	StartTime time.Time
	EndTime   time.Time
	Time      time.Duration
	Output    string
	Error     bool
}

var (
	commands = make(map[string]command)
)
var fallback func(string) (string, error)

func RegisterCommand(name string, description string, help string,
	callback func(string) (string, error), tabComplete func(string, int) string) {

	debug("Registering command: %v: %v", name, description)
	commands[name] = command{name, description, help, callback, tabComplete}
}

func RegisterFallbackCommand(fb func(string) (string, error)) {
	fallback = fb
}

// Is a command registered and valid
func isValidCommand(command string) bool {
	_, ok := commands[command]
	return ok
}

func getCommandFromInput(input string) (command, error) {
	// Split into each argument
	args := strings.Fields(input)
	var cmd command
	if len(args) == 0 {
		return cmd, errors.New("No command given")
	}

	cmd, ok := commands[args[0]]
	if !ok {
		return cmd, errors.New(fmt.Sprintf("Command '%v' not found. Try 'help'", args[0]))
	}

	return commands[args[0]], nil
}

// Take current user input, and expand tabs
// If a command is already in args[0], call that command's tabComplete
func TabComplete(partial string, tabcount int) string {
	// TODO: Implement this
	return partial
}

// Take a command, and call the appropriate command's callback
func ProcessInput(input string) (output CommandOutput, err error) {
	if len(commands) == 0 {
		return output, errors.New("No commands registered")
	}
	if len(input) == 0 {
		return output, nil
	}

	cmd, err := getCommandFromInput(input)
	output.StartTime = time.Now()

	if err == nil {
		args := strings.Split(input, " ")
		subcmd := strings.Join(args[1:], " ")
		// Call function with arguments
		output.Output, err = cmd.callback(subcmd)
	} else {
		// Invalid command
		//output.Output = fmt.Sprintf("%v", err)
		output.Output, err = fallback(input)
	}

	output.EndTime = time.Now()
	output.Time = output.EndTime.Sub(output.StartTime)

	return output, err
}
