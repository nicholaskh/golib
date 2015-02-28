package server

import (
	"net"
	"time"

	log "github.com/nicholaskh/log4go"
)

type Handler interface {
	Run(conn net.Conn)
}

func (this *Server) LaunchTcpServ(handler Handler) (err error) {
	ln, err := net.Listen("tcp", this.servAddr)

	if err != nil {
		log.Error("Launch tcp server error: %s", err.Error())
	}

	this.fd = ln

	log.Info("Listening on %s", this.servAddr)

	for {
		conn, err := this.fd.Accept()
		if err != nil {
			log.Error("Accept error: %s", err.Error())
		}

		handler.Run(conn)
	}
}

func (this *Server) StopTcpServ() {
	this.fd.Close()
	log.Info("HTTP server stopped")
}

func (this *Server) PingClient(conn net.Conn, interval time.Duration) {
	select {
	case <-time.Tick(interval):
		ping(conn)
	}
}
