package server

import "net"

type Protocol interface {
	Marshal([]byte) []byte
	SetConn(conn net.Conn)
	Read() ([]byte, error)
}
