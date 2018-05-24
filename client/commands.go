package client

import (
	"fmt"
	"github.com/xunil154/gobar/ui"
	"os"
	"os/exec"
	"strings"
)

func BootstrapCommands() {
	ui.RegisterCommand("chargen", "Generate characters to help with overflows",
		"Generates a set of strings that could aid in developing exploits for"+
			" buffer overflows",
		chargen, chargenTabComplete)
	ui.RegisterFallbackCommand(execFallback)
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
