package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	*Server
	router *mux.Router
}

func NewHttpServer(name string) *HttpServer {
	this := new(HttpServer)
	this.Server = NewServer(name)
	this.router = mux.NewRouter()

	return this
}

func (this *HttpServer) Launch(listenAddr string) {
	http.ListenAndServe(listenAddr, this.router)
}

func (this *HttpServer) RegisterHandler(path string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	return this.router.HandleFunc(path, handler)
}
