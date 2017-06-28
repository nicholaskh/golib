package server

import (
	"net"
	"sync"
	"errors"
	"bytes"
	"encoding/binary"

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
		return errors.New("link has been closed")
	}
	_, err := this.Conn.Write(data)
	this.Mutex.Unlock()

	if this.conn_type == CONN_TYPE_LONG_POLLING {
		this.Close()
	}
	return err
}

func (this *Client) WriteFormatMsg(op, body string) error {

	dataBuff := bytes.NewBuffer([]byte{})

	opBytes := []byte(op)
	var bodyBytes []byte
	if body != ""{
		bodyBytes = []byte(body)
	}

	// write op to dataBuff
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(len(opBytes)))
	dataBuff.Write(buf.Bytes())

	buf.Reset()
	binary.Write(buf, binary.BigEndian, opBytes)
	dataBuff.Write(buf.Bytes())

	// write body to dataBuff
	if body != "" {
		buf.Reset()
		binary.Write(buf, binary.BigEndian, int32(len(bodyBytes)))
		dataBuff.Write(buf.Bytes())

		buf.Reset()
		binary.Write(buf, binary.BigEndian, bodyBytes)
		dataBuff.Write(buf.Bytes())
	}

	// write into io
	return this.WriteBinMsg(dataBuff.Bytes())

}

func (this *Client) WriteBinMsg(msg []byte) error {
	data := this.Proto.Marshal(msg)
	this.Mutex.Lock()
	if !this.IsConnected() {
		return errors.New("link has been closed")
	}
	_, err := this.Conn.Write(data)
	this.Mutex.Unlock()

	if this.conn_type == CONN_TYPE_LONG_POLLING {
		this.Close()
	}

	return err
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
