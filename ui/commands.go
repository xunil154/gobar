package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func BootstrapCommands() {
	RegisterCommand("help", "Display help information", "Show this message",
		help, TabComplete)

	RegisterCommand("chargen", "Generate characters to help with overflows",
		"Generates a set of strings that could aid in developing exploits for"+
			" buffer overflows",
		chargen, chargenTabComplete)
	RegisterFallbackCommand(execFallback)
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

func execFallback(command string) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin   // Pass our stdin to cmd stdin
	cmd.Stdout = os.Stdout // Cmd stdout to ours
	err := cmd.Run()
	return "", err
}

func chargen(command string) (string, error) {
	pattern := "ABCDEF0123456789"
	ret := ""
	genby := func(s string, count int) (gen string) {
		for i := 0; i < count; i++ {
			gen += s
		}
		return gen
	}
	genbar := func(count int) string {
		ret := fmt.Sprintf("x%02d ", count)
		for i := 0; i*count <= 60; i++ {
			is := fmt.Sprintf("%d", i*count)
			ns := fmt.Sprintf("%d", (i+1)*count)
			buf := genby(" ", count-len(ns))

			ret += is + buf
		}
		ret += "\n"
		return ret
	}

	ranges := []int{10, 8, 5, 4}

	for _, inc := range ranges {
		ret += "\n" + genbar(inc)
		for i := inc; i <= 60; i += inc {
			gen := ""
			for j := 0; len(gen) < i; j++ {
				gen += genby(pattern[j:j+1], inc)
			}
			ret += fmt.Sprintf("%04d %v\n", len(gen), gen)
		}
	}

	return ret, nil
}

func chargenTabComplete(partial string, tabcount int) string {
	return ""
}
