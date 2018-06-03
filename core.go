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
	return ui.PromptSegment{"core", "black", "green"}
}

func registerCommands() {
	ui.BootstrapCommands()
	core.BootstrapCommands()
}

func main() {
	flag.Parse()
	registerCommands()
	ui.AddSegment(defaultPrompt())
	for {
		input := ui.GetUserInput(ui.TabComplete)
		if input == "exit" || input == "quit" {
			break
		}
		fmt.Println("")
		output, err := ui.ProcessInput(input)

		if err != nil {
			ui.Error(fmt.Sprintf("%v", err))
		} else if len(output.Output) > 0 {
			ui.Output(output.Output)
		}
	}

	core.Shutdown()
	ui.Exit()
}
