package server

import (
	"net"
	"time"

	"github.com/nicholaskh/golib/sync2"
	log "github.com/nicholaskh/log4go"
)

type TcpServer struct {
	*Server
	SessTimeout         time.Duration
	clientProcessor     ClientProcessor
	AcceptLock          *sync2.Semaphore
	initialGoRoutineNum int
}

type Client struct {
	net.Conn
	LastTime time.Time
	Ticker   *time.Ticker
	Done     chan byte
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
	this.AcceptLock = sync2.NewSemaphore(1, 0)
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

func (this *TcpServer) CheckTimeout(client *Client, close func() error) {
	for {
		select {
		case <-client.Ticker.C:
			log.Debug("Check client timeout: %s", client.Conn.RemoteAddr())
			if time.Now().After(client.LastTime.Add(this.SessTimeout)) {
				log.Warn("Client connection timeout: %s", client.Conn.RemoteAddr())
				close()
				return
			}

		case <-client.Done:
			close()
			return
		}
	}
}
