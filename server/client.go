package server

import (
	"net"
	"sync"
	"time"

	log "github.com/nicholaskh/log4go"
)

const (
	CONN_TYPE_TCP = iota
	CONN_TYPE_LONG_POLLING
)

type Client struct {
	net.Conn
	sync.Mutex
	conn_type int8
}

func NewClient(conn net.Conn, now time.Time, ctype int8) *Client {
	return &Client{Conn: conn, conn_type: ctype}
}

func (this *Client) WriteMsg(msg string) {
	this.Conn.Write([]byte(msg))
	if this.conn_type == CONN_TYPE_LONG_POLLING {
		this.Close()
	}
}

func (this *Client) IsConnected() bool {
	return this.Conn != nil
}

// reentrant safe
func (this *Client) Close() {
	if this.Conn == nil {
		return
	}
	this.Mutex.Lock()
	log.Info("Client shutdown: %s", this.Conn.RemoteAddr())
	err := this.Conn.Close()
	if err != nil {
		log.Error(err)
	}
	this.Conn = nil
	this.Mutex.Unlock()
}
