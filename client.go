package main

import (
	"fmt"
	"github.com/vharitonsky/iniflags"
	"github.com/xunil154/gobar/client"
	"github.com/xunil154/gobar/ui"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Ltime)
)

func defaultPrompt() ui.PromptSegment {
	return ui.PromptSegment{"gobar", "black", "white"}
}

func registerCommands() {
}

func main() {
	iniflags.Parse()

	ui.BootstrapCommands()
	client.BootstrapCommands()
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

	ui.Exit()
}
