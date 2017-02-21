package server

import "net"

const (
	ProtoType = iota
	SIMPLE
	FIXED_LENGTH
)

type Protocol interface {
	Marshal([]byte) []byte
	SetConn(conn net.Conn)
	Read() ([]byte, error)
}

func factoryProto(protoType int) (proto Protocol) {
	switch protoType {
	case SIMPLE:
		proto = NewSimpleProtocol()
	case FIXED_LENGTH:
		proto = NewFixedLengthProtocol()
	default:
		proto = NewFixedLengthProtocol()
	}
	return proto
}
