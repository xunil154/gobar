package main

import (
	"fmt"
	"github.com/xunil154/gobar/ui"
	"log"
	"os"
)

var (
	uiSegments = make([]ui.PromptSegment, 0, 10)
	logger     = log.New(os.Stdout, "logger: ", log.Ltime)
)

func defaultPrompt() ui.PromptSegment {
	return ui.NewPromptSegment("gobar", "black", "white")
}

func registerCommands() {
}

func main() {
	// Shared with commands
	uiSegments = append(uiSegments, defaultPrompt())

	ui.BootstrapCommands()
	for {
		input := ui.GetUserInput(uiSegments, ui.TabComplete)
		if input == "exit" || input == "quit" {
			break
		}
		fmt.Println("")
		output, err := ui.ProcessInput(input)

		if err != nil {
			ui.Error(fmt.Sprintf("%v", err), uiSegments)
		} else if len(output.Output) > 0 {
			ui.Output(output.Output, uiSegments)
		}
	}

	ui.Exit()
}

////// COMMANDS \\\\\\\
