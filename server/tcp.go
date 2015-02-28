package server

import (
	"fmt"
	"net"
)

type Handler interface {
	run(conn net.Conn)
}

func LaunchTcpServ(listenAddr string, handler Handler) {
	ln, err := net.Listen("tcp", listenAddr)

	// TODO log
	if err != nil {

	}

	defer ln.Close()

	// TODO log
	fmt.Println("Listening on " + listenAddr)

	for {
		conn, err := ln.Accept()

		// TODO log
		if err != nil {

		}

		handler.run(conn)
	}
}
