package server

import (
	"net"
	"time"

	log "github.com/nicholaskh/log4go"
)

type Handler interface {
	Run(*Client)
}

type Client struct {
	net.Conn
	LastTime time.Time
}

func (this *Server) LaunchTcpServ(listenAddr string, handler Handler, servTimeout time.Duration) (err error) {
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

		// TODO use a thread pool
		client := &Client{Conn: conn, LastTime: time.Now()}

		go handler.Run(client)
		if servTimeout.Nanoseconds() > int64(0) {
			go this.checkTimeout(client, servTimeout)
		}
	}
}

func (this *Server) StopTcpServ() {
	this.fd.Close()
	log.Info("HTTP server stopped")
}

func (this *Server) checkTimeout(client *Client, timeout time.Duration) {
	for {
		select {
		case <-time.Tick(timeout):
			log.Debug("Check client timeout: %s", client.Conn.RemoteAddr())
			if time.Now().After(client.LastTime.Add(timeout)) {
				log.Warn("Client connection timeout: %s", client.Conn.RemoteAddr())
				client.Conn.Close()
				return
			}
		}
	}
}
