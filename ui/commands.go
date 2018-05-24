package ui

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

func BootstrapCommands() {
	RegisterCommand("help", "Display help information", "Show this message",
		help, TabComplete)

	RegisterFallbackCommand(help) // Display help if unknown
}

func defaultHelp() string {
	var help string = "Available commands:\n"
	cmds := make([]string, 0, len(commands))
	for cmd, _ := range commands {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	for _, cmd := range cmds {
		help += fmt.Sprintf("\t%v - %v\n", cmd, commands[cmd].help)
	}
	help += "\texit - exit the application"
	return help
}

func help(command string) (string, error) {
	args := strings.Fields(command)
	if len(args) > 0 && !isValidCommand(args[0]) {
		args := strings.Split(command, " ")
		return "", errors.New(fmt.Sprintf("Command '%v' not found", args[0]))
	} else if len(args) > 0 {
		return commands[args[0]].help, nil
	}
	return defaultHelp(), nil
}
