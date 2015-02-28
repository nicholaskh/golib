package server

import (
	"net"
	"time"

	log "github.com/nicholaskh/log4go"
)

type Handler interface {
	Run(conn net.Conn)
}

func (this *Server) LaunchTcpServ(listenAddr string, handler Handler, pingInterval time.Duration) (err error) {
	ln, err := net.Listen("tcp", listenAddr)

	if err != nil {
		log.Error("Launch tcp server error: %s", err.Error())
	}

	this.fd = ln

	log.Info("Listening on %s", listenAddr)

	for {
		conn, err := this.fd.Accept()
		if err != nil {
			log.Error("Accept error: %s", err.Error())
		}

		handler.Run(conn)
		if pingInterval.Nanoseconds() > int64(0) {
			go this.PingClient(conn, pingInterval)
		}
	}
}

func (this *Server) StopTcpServ() {
	this.fd.Close()
	log.Info("HTTP server stopped")
}

// TODO retry
func (this *Server) PingClient(conn net.Conn, interval time.Duration) {
	for {
		select {
		case <-time.Tick(interval):
			log.Debug("Ping client %s", conn.RemoteAddr())
			_, err := conn.Write([]byte{0})
			if err != nil {
				conn.Close()
				return
			}
		}
	}
}
