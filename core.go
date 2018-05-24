package main

import (
	"flag"
	"fmt"
	"github.com/xunil154/gobar/core"
	"github.com/xunil154/gobar/ui"
	"log"
	"os"
)

var (
	uiSegments = make([]ui.PromptSegment, 0, 10)
	logger     = log.New(os.Stdout, "logger: ", log.Ltime)
)

func defaultPrompt() ui.PromptSegment {
	return ui.PromptSegment{"gobar", "black", "green"}
}

func registerCommands() {
	ui.BootstrapCommands()
	core.BootstrapCommands()
}

func main() {
	flag.Parse()
	// Shared with commands
	uiSegments = append(uiSegments, defaultPrompt())

	registerCommands()

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
