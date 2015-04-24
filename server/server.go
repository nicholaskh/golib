package server

import (
	"net"
	"os"
	"syscall"
	"time"

	"github.com/nicholaskh/golib/signal"
)

type Server struct {
	Name       string
	configFile string
	StartedAt  time.Time
	pid        int
	hostname   string
	Fd         net.Listener
}

func NewServer(name string) (this *Server) {
	this = new(Server)
	this.Name = name

	this.StartedAt = time.Now()
	this.hostname, _ = os.Hostname()
	this.pid = os.Getpid()
	signal.IgnoreSignal(syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGSTOP)

	return
}
