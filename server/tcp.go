package server

import (
	"net"
	"time"

	"github.com/nicholaskh/golib/sync2"
	log "github.com/nicholaskh/log4go"
)

type Handler interface {
	Run(*Client)
}

type TcpServer struct {
	*Server
	sessTimeout         time.Duration
	handler             Handler
	acceptLock          *sync2.Semaphore
	initialGoRoutineNum int
}

type Client struct {
	net.Conn
	LastTime time.Time
	ticker   *time.Ticker
	done     chan byte
}

func NewTcpServer(name string) (this *TcpServer) {
	this = new(TcpServer)
	this.Server = NewServer(name)

	return
}

func (this *TcpServer) LaunchTcpServ(listenAddr string, handler Handler, sessTimeout time.Duration, initialGoRoutineNum int) (err error) {
	this.sessTimeout = sessTimeout
	this.handler = handler
	this.acceptLock = sync2.NewSemaphore(1, 0)
	this.initialGoRoutineNum = initialGoRoutineNum
	tcpAddr, _ := net.ResolveTCPAddr("tcp", listenAddr)
	ln, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Error("Launch tcp server error: %s", err.Error())
	}

	this.fd = ln

	log.Info("Listening on %s", listenAddr)

	for i := 0; i < int(this.initialGoRoutineNum); i++ {
		go this.startGoRoutine()
	}

	return
}

func (this *TcpServer) startGoRoutine() {
	log.Debug("start server go routine")
	this.acceptLock.Acquire()
	conn, err := this.fd.(*net.TCPListener).AcceptTCP()
	this.acceptLock.Release()

	go this.startGoRoutine()

	if err != nil {
		log.Error("Accept error: %s", err.Error())
	}

	client := &Client{Conn: conn, LastTime: time.Now(), ticker: time.NewTicker(this.sessTimeout), done: make(chan byte)}

	if this.sessTimeout.Nanoseconds() > int64(0) {
		go this.checkTimeout(client)
	}
	this.handler.Run(client)
	client.done <- 0

}

func (this *TcpServer) StopTcpServ() {
	this.fd.Close()
	log.Info("HTTP server stopped")
}

func (this *TcpServer) checkTimeout(client *Client) {
	for {
		select {
		case <-client.ticker.C:
			log.Debug("Check client timeout: %s", client.Conn.RemoteAddr())
			if time.Now().After(client.LastTime.Add(this.sessTimeout)) {
				log.Warn("Client connection timeout: %s", client.Conn.RemoteAddr())
				client.Conn.Close()
				return
			}

		case <-client.done:
			return
		}
	}
}
