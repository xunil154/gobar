package core

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func listenLoop(proto string, port int, handler func(net.Conn)) error {
	if proto != "tcp" && proto != "udp" {
		return errors.New(fmt.Sprintf("Invalid protocol specified %v must be either 'tcp' or 'udp'", proto))
	}

	listener, err := net.Listen(proto, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()
		go handler(conn)
	}

	return nil
}

func defaultConnectionHandler(connection net.Conn) {
	log.Print("New connection accepted from %v", connection.RemoteAddr())
	io.Copy(connection, os.Stdout) // print all output for now
}
