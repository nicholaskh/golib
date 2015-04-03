package server

import (
	conf "github.com/nicholaskh/jsconf"
	log "github.com/nicholaskh/log4go"
)

// universal config keys:
// max_cpu
func LoadConfig(fn string) *conf.Conf {
	log.Info("Loading config file: %s", fn)

	var err error
	config, err := conf.Load(fn)
	if err != nil {
		panic(err)
	}

	return config
}
