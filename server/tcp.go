package server

import (
	"net"
	"sync"
	"time"

	log "github.com/nicholaskh/log4go"
)

type TcpServer struct {
	*Server
	SessTimeout         time.Duration
	clientProcessor     ClientProcessor
	AcceptLock          sync.Mutex
	initialGoRoutineNum int
}

type Client struct {
	net.Conn
	LastTime    time.Time
	sessTimeout time.Duration
	Done        chan byte
	Closed      bool
	sync.Mutex
	OnClose func()
}

func NewClient(conn net.Conn, now time.Time, sessTimeout time.Duration) *Client {
	return &Client{Conn: conn, LastTime: now, sessTimeout: sessTimeout, Done: make(chan byte), Closed: false}
}

type ClientProcessor interface {
	Run()
}

func (this *Client) WriteMsg(msg string) {
	this.Conn.Write([]byte(msg))
}

func NewTcpServer(name string) (this *TcpServer) {
	this = new(TcpServer)
	this.Server = NewServer(name)

	return
}

func (this *TcpServer) LaunchTcpServer(listenAddr string, clientProcessor ClientProcessor, sessTimeout time.Duration, initialGoRoutineNum int) (err error) {
	this.SessTimeout = sessTimeout
	this.clientProcessor = clientProcessor
	this.initialGoRoutineNum = initialGoRoutineNum
	tcpAddr, _ := net.ResolveTCPAddr("tcp", listenAddr)
	ln, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Error("Launch tcp server error: %s", err.Error())
	}

	this.Fd = ln

	log.Info("Listening on %s", listenAddr)

	for i := 0; i < int(this.initialGoRoutineNum); i++ {
		go this.clientProcessor.Run()
	}

	return
}

func (this *TcpServer) StopTcpServ() {
	this.Fd.Close()
	log.Info("HTTP server stopped")
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
