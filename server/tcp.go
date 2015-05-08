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

type ClientProcessor interface {
	OnAccept(*Client)
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
		go this.startProcessorThread()
	}

	return
}

func (this *TcpServer) startProcessorThread() {
	log.Debug("start server go routine")
	this.AcceptLock.Lock()
	conn, err := this.Fd.(*net.TCPListener).AcceptTCP()
	this.AcceptLock.Unlock()

	if this.Fd == nil {
		return
	}

	if err != nil {
		log.Error("Accept error: %s", err.Error())
	}

	go this.startProcessorThread()
	if conn == nil {
		return
	}
	client := NewClient(conn, time.Now(), CONN_TYPE_TCP)
	this.clientProcessor.OnAccept(client)
}

func (this *TcpServer) StopTcpServer() {
	this.Fd.Close()
	this.Fd = nil
	log.Info("TCP server[%s] stopped", this.Name)
}
