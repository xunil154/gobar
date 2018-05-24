package core

import (
	"errors"
	"fmt"
	"github.com/xunil154/gobar/ui"
	"strconv"
	"strings"
)

var (
	listenPort = 8443
	listenHelp = "listen <port>\nListen for incomming connections on defined port"
)

func BootstrapCommands() {
	ui.RegisterCommand("listen", "Listen for incomming connections",
		listenHelp, listen, listenTabComplete)
}

func listen(command string) (string, error) {
	proto := "tcp"

	args := strings.Fields(command)
	if len(args) != 1 {
		return "", errors.New("No port specified\n" + listenHelp)
	}
	port, err := strconv.Atoi(args[0])
	if err != nil || port < 0 || port > 65535 {
		return "", errors.New("Invalid port")
	}

	go listenLoop(proto, port, defaultConnectionHandler)

	return fmt.Sprintf("Listening on %v/%v", port, proto), nil
}

func listenTabComplete(partial string, tabcount int) string {
	ui.Debug("Listen TC: %v", tabcount)
	if tabcount == 1 {
		return "TCP_PORT"
	}
	return ""
}
