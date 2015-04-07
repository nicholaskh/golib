package server

import (
	"net"
	"sync"
	"time"

	log "github.com/nicholaskh/log4go"
)

const (
	CTYPE_TCP = iota
	CTYPE_LONG_POLLING
)

type Client struct {
	net.Conn
	LastTime    time.Time
	sessTimeout time.Duration
	Done        chan byte
	Closed      bool
	sync.Mutex
	OnClose func()
	ctype   int8
}

func NewClient(conn net.Conn, now time.Time, sessTimeout time.Duration, ctype int8) *Client {
	return &Client{Conn: conn, LastTime: now, sessTimeout: sessTimeout, Done: make(chan byte), Closed: false, ctype: ctype}
}

func (this *Client) WriteMsg(msg string) {
	this.Conn.Write([]byte(msg))
	if this.ctype == CTYPE_LONG_POLLING {
		this.Close()
	}
}

func (this *Client) CheckTimeout() {
	ticker := time.NewTicker(this.sessTimeout)
	for {
		select {
		case <-ticker.C:
			log.Debug("Check client timeout: %s", this.Conn.RemoteAddr())
			if time.Now().After(this.LastTime.Add(this.sessTimeout)) {
				log.Warn("Client connection timeout: %s", this.Conn.RemoteAddr())
				this.Close()
				return
			}

		case <-this.Done:
			this.Close()
			return
		}
	}
}

func (this *Client) Close() {
	if this.OnClose != nil {
		this.OnClose()
	}
	this.Mutex.Lock()
	this.Closed = true
	err := this.Conn.Close()
	if err != nil {
		log.Error(err)
	}
	this.Mutex.Unlock()
}
