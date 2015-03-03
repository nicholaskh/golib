package slab

import (
	"github.com/nicholaskh/pretty"
	"log"
)

func debug(format string, args ...interface{}) {
	const debug = true
	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println(pretty.Sprintf(format, args...))
	}
}
