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
	net.Conn
}

func NewClient(conn net.Conn, now time.Time, ctype int8, proto *Protocol) *Client {
	return &Client{conn_type: ctype, Proto: proto, Conn: conn}
}

func (this *Client) WriteMsg(msg string) {
	data := this.Proto.Marshal([]byte(msg))
	this.Conn.Write(data)
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
