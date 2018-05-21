package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func BootstrapCommands() {
	RegisterCommand("help", "Display help information", "Show this message",
		help, TabComplete)
	RegisterFallbackCommand(execFallback)
}

func defaultHelp() string {
	var help string = "Available commands:\n"
	for cmd, command := range commands {
		help += fmt.Sprintf("\t%v - %v\n", cmd, command.help)
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

func execFallback(command string) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin   // Pass our stdin to cmd stdin
	cmd.Stdout = os.Stdout // Cmd stdout to ours
	err := cmd.Run()
	return "", err
}
