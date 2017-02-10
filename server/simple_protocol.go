package server

import (
	"net"
)

type SimpleProtocol struct {
	net.Conn
}

func NewSimpleProtocol() *SimpleProtocol {
	this := new(SimpleProtocol)
	return this
}

func (this *SimpleProtocol) Marshal(payload []byte) []byte {
	return payload
}

func (this *SimpleProtocol) SetConn(conn net.Conn) {
	this.Conn = conn
}

func (this *SimpleProtocol) Read() (buff []byte, err error) {
	buff = make([]byte, 1024)
	_, err = this.Conn.Read(buff)
	return
}
