package server

import (
	conf "github.com/nicholaskh/jsconf"
	log "github.com/nicholaskh/log4go"
)

// universal config keys:
// max_cpu
func (this *Server) LoadConfig(fn string) *Server {
	log.Info("Server[%s %s@%s] loading config file: %s", this.Name, BuildID, VERSION, fn)
	this.configFile = fn

	var err error
	this.Conf, err = conf.Load(fn)
	if err != nil {
		panic(err)
	}

	this.servAddr = this.Conf.String("serv_addr", ":2222")

	return this
}
