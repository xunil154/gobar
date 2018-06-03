package core

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Listener interface {
	Listen(func(Handler))
	Stop()
	String() string
}

type Handler interface {
	Handle(func(Handler)) // Function to call when disconnected
	String() string
	Stop()
}

// Internal Agents
type agentHandler struct {
	listener Listener // parent listener
	client   *net.TCPConn
	stop     chan int
}

type agentListener struct {
	port     int
	protocol string
	listener *net.TCPListener
	stop     chan int
}

func NewListener(protocol string, port int) (Listener, error) {
	var x Listener
	if protocol != "tcp" {
		return x, errors.New("Only TCP is supported for now :(")
	}

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{nil, port, "::"})
	if err != nil {
		var x Listener
		return x, err
	}

	return agentListener{
		port,
		protocol,
		listener,
		make(chan int), // Stop channel
	}, nil
}

func NewHandler(client *net.TCPConn, listener Listener) Handler {
	return agentHandler{
		listener,
		client,
		make(chan int), // Stop channel
	}
}

func (listener agentListener) String() string {
	return fmt.Sprintf("%v/%v", listener.port, listener.protocol)
}

func (listener agentListener) Stop() {
	fmt.Printf("Shutting down agent: %v", listener)
	listener.listener.Close()
	listener.stop <- 1
}
func (listener agentListener) Listen(callback func(Handler)) {
	newConnections := make(chan *net.TCPConn)

	// This will send new connections through a channel so our switch case works
	go func(l *net.TCPListener) {
		// Constantly accept new connections
		for {
			l.SetDeadline(time.Now().Add(time.Duration(50) * time.Millisecond))

			c, err := l.AcceptTCP()
			if err == nil {
				newConnections <- c
			}
			// TODO: Better error handeling
			//newConnections <- nil
		}
	}(listener.listener)

	// Wait for connections or call to stop
	for {
		select {
		case client := <-newConnections:
			fmt.Println("Accepted new connection")
			agent := NewHandler(client, listener)
			callback(agent)
		case stop := <-listener.stop:
			if stop == 1 {
				fmt.Println("Shutting down listener %s/%v", listener.protocol, listener.port)
				break
			}
		default:
			time.Sleep(time.Duration(10) * time.Millisecond) // sleep and try again
		}
	}
	// Shutdown stuff
	return
}

func (handler agentHandler) String() string {
	return fmt.Sprintf("%v - %v", handler.listener, handler.client.RemoteAddr())
}
func (handler agentHandler) Handle(shutdown func(Handler)) {
	handler.client.Write([]byte("Greeting!\n"))
	fmt.Printf("Handeling connection for listener: %v\n", handler.listener)
	time.Sleep(time.Duration(4) * time.Second)
	handler.client.Close() // TODO: Not do this
	shutdown(handler)
}
func (handler agentHandler) Stop() {
	fmt.Printf("Shutting down handler: %v\n", handler)
	handler.stop <- 1
}
