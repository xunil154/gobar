package core

import (
	"errors"
	"fmt"
	"github.com/xunil154/gobar/ui"
	"sort"
	"strconv"
	"strings"
)

var (
	listenPort   = 8443
	listenHelp   = "listen <port>\nListen for incomming connections on defined port"
	nolistenHelp = "nolisten <port>\nStop listening on given port"

	listeners = make(map[int]Listener) // id -> listener
	agents    = make(map[int]Handler)  // id -> agent
)

// Return an ID for a given map
func nextAgentId(x map[int]Handler) int {
	// Look for the first available
	for i := 0; i < len(x); i++ {
		_, ok := x[i]
		if !ok {
			return i
		}
	}
	// If none exist, it is the size of the map
	return len(x)
}
func nextListenerId(x map[int]Listener) int {
	// Look for the first available
	for i := 0; i < len(x); i++ {
		_, ok := x[i]
		if ok {
			return i
		}
	}
	// If none exist, it is the size of the map
	return len(x)
}

func acceptHandler(handler Handler) {
	id := nextAgentId(agents)
	agents[id] = handler
	fmt.Printf("\nNew Agent: %d - %v\n", id, handler.String())

	go handler.Handle(disconnectHandler)
}
func disconnectHandler(handler Handler) {
	for id, agent := range agents {
		if agent == handler { // **should** work https://golang.org/ref/spec#Comparison_operators
			fmt.Printf("\nAgent disconnect: [%d] %v\n", id, handler)
			delete(agents, id)
			return
		}
	}
	fmt.Printf("\nUnknown agent disconnect [?] %v\n", handler)
}

func BootstrapCommands() {
	ui.RegisterCommand("nolisten", "Stop listening for incomming connections",
		nolistenHelp, noListen, noListenTabComplete)
	ui.RegisterCommand("listen", "Listen for incomming connections",
		listenHelp, listen, listenTabComplete)
	ui.RegisterCommand("agents", "List connected agents",
		"Lists connected agents", listAgents, nil)
}

func listAgents(command string) (string, error) {
	output := "\nID\tAgent\n"

	ids := make([]int, 0, len(agents))
	for k := range agents {
		ids = append(ids, k)
	}
	sort.Ints(ids)

	for id := range ids {
		agent, ok := agents[id]
		if !ok {
			fmt.Printf("[E] Bad agent id: %d", id)
		} else {
			output += fmt.Sprintf("%d\t%v\n", id, agent.String())
		}
	}
	return output, nil
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

	listener, err := NewListener(proto, port)
	if err != nil {
		return "", err
	}
	id := nextListenerId(listeners)
	listeners[id] = listener

	ui.AddSegment(ui.NewPromptSegment(fmt.Sprintf("🖧 %v", port), "black", "blue"))

	// Start the listener thread
	// Can be shut down by calling listener.Stop()
	go listener.Listen(acceptHandler)

	return fmt.Sprintf("Listening on %v/%v", port, proto), nil
}

func listenTabComplete(partial string, tabcount int) string {
	if tabcount == 1 {
		return "TCP_PORT"
	}
	return ""
}

func noListen(command string) (response string, err error) {
	for _, listener := range listeners {
		if command == listener.String() {
			response = fmt.Sprintf("Shutting down listener %v", listener)
			listener.Stop()
			break
		}
	}
	if response == "" {
		return "", errors.New(fmt.Sprintf("No listener found '%v'", command))

	}
	return response, nil
}
func noListenTabComplete(partial string, tabcount int) string {
	if tabcount == 0 {
		if len(listeners) == 1 {
			return listeners[0].String()
		}
	}
	if tabcount == 1 {
		options := make([]string, 0)
		for _, option := range listeners {
			options = append(options, option.String())
		}
		return strings.Join(options, "\t")
	}
	return ""
}

func Shutdown() {
	fmt.Println("\nShutting down agents/listeners\n")
	for _, agent := range agents {
		agent.Stop()
	}
	for _, listener := range listeners {
		listener.Stop()
	}
}
