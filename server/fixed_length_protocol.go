package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"net"

	log "github.com/nicholaskh/log4go"
)

const (
	HEAD_LENGTH = 4
)

type FixedLengthProtocol struct {
	net.Conn
}

func NewFixedLengthProtocol() *FixedLengthProtocol {
	this := new(FixedLengthProtocol)

	return this
}

func (this *FixedLengthProtocol) SetConn(conn net.Conn) {
	this.Conn = conn
}

//len+payload
func (this *FixedLengthProtocol) Marshal(payload []byte) []byte {
	buf := bytes.NewBuffer([]byte{})
	dataBuff := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(len(payload)))
	binary.Write(dataBuff, binary.BigEndian, payload)
	buf.Write(dataBuff.Bytes())

	return buf.Bytes()
}

func (this *FixedLengthProtocol) Read() ([]byte, error) {
	buf := make([]byte, HEAD_LENGTH)
	err := this.ReadN(this.Conn, buf, HEAD_LENGTH)
	if err != nil {
		if err != io.EOF {
			log.Error("[Protocol] Read data length error: %s", err.Error())
		}
		return []byte{}, err
	}
	//data length
	b_buf := bytes.NewBuffer(buf[:4])
	var dataLength int32
	binary.Read(b_buf, binary.BigEndian, &dataLength)

	//app + data
	payloadLength := int(dataLength)
	if payloadLength > math.MaxInt64 || payloadLength < 0 {
		return []byte{}, errors.New("[Protocol] Payload out of length")
	}
	payload := make([]byte, payloadLength)
	err = this.ReadN(this.Conn, payload, payloadLength)
	if err != nil && err != io.EOF {
		log.Error("[Protocol] Read data error: %s", err.Error())
		return []byte{}, err
	}

	return payload, nil
}

func (this *FixedLengthProtocol) ReadN(conn net.Conn, buf []byte, n int) (err error) {
	buffer := bytes.NewBuffer([]byte{})
	for n > 0 {
        var readN int
		b_buf := make([]byte, n)
		readN, err = conn.Read(b_buf)
		if err != nil {
			return err
		}
		n -= readN
		buffer.Write(b_buf)
	}
    _, err = buffer.Read(buf)
	return err
}
