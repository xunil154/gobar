package main

import (
	"fmt"
	"github.com/vharitonsky/iniflags"
	"github.com/xunil154/gobar/ui"
	"log"
	"os"
)

var (
	uiSegments = make([]ui.PromptSegment, 0, 10)
	logger     = log.New(os.Stdout, "logger: ", log.Ltime)
)

func defaultPrompt() ui.PromptSegment {
	return ui.PromptSegment{"gobar", "black", "white"}
}

func main() {
	iniflags.Parse()
	// Shared with commands
	uiSegments = append(uiSegments, defaultPrompt())

	for {
		input := ui.GetUserInput(uiSegments, ui.TabComplete)
		if input == "exit" || input == "quit" {
			break
		}
		fmt.Println("")
		output, err := ui.ProcessInput(input, &uiSegments)

		if err != nil {
			fmt.Println("[E]", err)
		}
		if len(output.Output) > 0 {
			fmt.Println(output.Output)
		}
	}

	ui.Exit()
}
