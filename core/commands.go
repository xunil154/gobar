package core

import (
	"github.com/xunil154/gobar/ui"
)

func BootstrapCommands() {
	ui.RegisterCommand("listen", "Listen for incomming connections",
		"listen <port>\nListen for incomming connections on defined port",
		listen, listenTabComplete)
}

func listen(command string) (string, error) {
	return "NOT IMPLEMENTED", nil
}

func listenTabComplete(partial string, tabcount int) string {
	ui.Debug("Listen TC: %v", tabcount)
	if tabcount == 1 {
		return "TCP_PORT"
	}
	return ""
}
