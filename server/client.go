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
	sync.Mutex
	conn_type int8
	Proto     *Protocol
}

func NewClient(conn net.Conn, now time.Time, ctype int8, proto *Protocol) *Client {
	return &Client{conn_type: ctype, Proto: proto}
}

func (this *Client) WriteMsg(msg string) {
	this.Proto.Conn.Write([]byte(msg))
	if this.conn_type == CONN_TYPE_LONG_POLLING {
		this.Close()
	}
}

func (this *Client) IsConnected() bool {
	return this.Proto.Conn != nil
}

// reentrant safe
func (this *Client) Close() {
	if this.Proto.Conn == nil {
		return
	}
	this.Mutex.Lock()
	log.Info("Client shutdown: %s", this.Proto.Conn.RemoteAddr())
	err := this.Proto.Conn.Close()
	if err != nil {
		log.Error(err)
	}
	this.Proto.Conn = nil
	this.Mutex.Unlock()
}
