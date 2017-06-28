package server

import (
	"net"
	"sync"

	log "github.com/nicholaskh/log4go"
)

const (
	CONN_TYPE_TCP = iota
	CONN_TYPE_LONG_POLLING
)

type Client struct {
	sync.Mutex
	conn_type int8
	Proto     Protocol
	net.Conn
}

func NewClient(conn net.Conn, ctype int8, proto Protocol) *Client {
	return &Client{conn_type: ctype, Proto: proto, Conn: conn}
}

func (this *Client) WriteMsg(msg string) error {
	data := this.Proto.Marshal([]byte(msg))
	this.Mutex.Lock()
	if !this.IsConnected() {
		return nil
	}
	_, err := this.Conn.Write(data)
	this.Mutex.Unlock()

	if this.conn_type == CONN_TYPE_LONG_POLLING {
		this.Close()
	}
	return err
}

func (this *Client) WriteBinMsg(msg []byte) {
	data := this.Proto.Marshal(msg)
	this.Mutex.Lock()
	if !this.IsConnected() {
		return
	}
	this.Conn.Write(data)
	this.Mutex.Unlock()

	if this.conn_type == CONN_TYPE_LONG_POLLING {
		this.Close()
	}
}

func (this *Client) IsConnected() bool {
	return this.Conn != nil
}

// reentrant safe
func (this *Client) Close() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

    if !this.IsConnected() {
		return
	}
	log.Info("Client shutdown: %s", this.Conn.RemoteAddr())
	err := this.Conn.Close()
	if err != nil {
		log.Error(err)
	}
	this.Conn = nil
}
