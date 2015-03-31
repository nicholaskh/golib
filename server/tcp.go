package server

import (
	"io"
	"net"
	"time"

	"github.com/nicholaskh/golib/sync2"
	log "github.com/nicholaskh/log4go"
)

type TcpServer struct {
	*Server
	sessTimeout         time.Duration
	clientHandler       ClientHandler
	acceptLock          *sync2.Semaphore
	initialGoRoutineNum int
}

type ClientHandler interface {
	OnAccept(*Client)
	OnRead(string)
	OnClose()
}

type Client struct {
	net.Conn
	LastTime time.Time
	ticker   *time.Ticker
	done     chan byte
}

func (this *Client) WriteMsg(msg string) {
	this.Conn.Write([]byte(msg))
}

func NewTcpServer(name string) (this *TcpServer) {
	this = new(TcpServer)
	this.Server = NewServer(name)

	return
}

func (this *TcpServer) LaunchTcpServer(listenAddr string, clientHandler ClientHandler, sessTimeout time.Duration, initialGoRoutineNum int) (err error) {
	this.sessTimeout = sessTimeout
	this.clientHandler = clientHandler
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
		go this.dealSingleClient()
	}

	return
}

func (this *TcpServer) dealSingleClient() {
	log.Debug("start server go routine")
	this.acceptLock.Acquire()
	conn, err := this.fd.(*net.TCPListener).AcceptTCP()
	this.acceptLock.Release()
	if err != nil {
		log.Error("Accept error: %s", err.Error())
	}

	go this.dealSingleClient()

	client := &Client{Conn: conn, LastTime: time.Now(), ticker: time.NewTicker(this.sessTimeout), done: make(chan byte)}
	this.clientHandler.OnAccept(client)

	if this.sessTimeout.Nanoseconds() > int64(0) {
		go this.checkTimeout(client)
	}

	for {
		input := make([]byte, 1460)
		n, err := client.Conn.Read(input)

		input = input[:n]

		if err != nil {
			if err == io.EOF {
				log.Info("Client shutdown: %s", client.Conn.RemoteAddr())
				this.clientHandler.OnClose()
				return
			} else if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				log.Error("Read from client[%s] error: %s", client.Conn.RemoteAddr(), err.Error())
				this.clientHandler.OnClose()
				return
			}
		}

		client.LastTime = time.Now()

		strInput := string(input)
		log.Debug("input: %s", strInput)

		this.clientHandler.OnRead(strInput)
	}

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
				this.clientHandler.OnClose()
				return
			}

		case <-client.done:
			this.clientHandler.OnClose()
			return
		}
	}
}
