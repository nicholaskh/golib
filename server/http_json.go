package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/nicholaskh/log4go"
	"io"
	"net"
	"net/http"
	"time"

	_ "expvar"         // localhost:xx/debug/vars
	_ "net/http/pprof" // localhost:xx/debug/pprof
)

type HttpJsonServer struct {
	*Server
	api *httpRestApi
}

type httpRestApi struct {
	httpListener net.Listener
	httpServer   *http.Server
	httpRouter   *mux.Router
	httpPaths    []string
}

func NewHttpJsonServer() *HttpJsonServer {
	this := new(HttpJsonServer)

	return this
}

func (this *HttpJsonServer) LaunchHttpServer(listenAddr string, debugAddr string) (err error) {
	if this.api != nil {
		return nil
	}

	this.api = new(httpRestApi)
	this.api.httpPaths = make([]string, 0, 10)
	this.api.httpRouter = mux.NewRouter()
	this.api.httpServer = &http.Server{Addr: listenAddr, Handler: this.api.httpRouter}

	this.api.httpListener, err = net.Listen("tcp", this.api.httpServer.Addr)
	if err != nil {
		this.api = nil
		return err
	}

	if debugAddr != "" {
		log.Debug("HTTP serving at %s with pprof at %s", listenAddr, debugAddr)
	} else {
		log.Debug("HTTP serving at %s", listenAddr)
	}

	go this.api.httpServer.Serve(this.api.httpListener)
	if debugAddr != "" {
		go http.ListenAndServe(debugAddr, nil)
	}

	return nil
}

func (this *HttpJsonServer) StopHttpServer() {
	if this.api != nil && this.api.httpListener != nil {
		this.api.httpListener.Close()
		this.api.httpListener = nil

		log.Info("HTTP server stopped")
	}
}

func (this *HttpJsonServer) Launched() bool {
	return this.api != nil
}

func (this *HttpJsonServer) RegisterHttpApi(path string,
	handlerFunc func(http.ResponseWriter,
		*http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	wrappedFunc := func(w http.ResponseWriter, req *http.Request) {
		var (
			ret interface{}
			t1  = time.Now()
		)

		params, err := this.api.decodeHttpParams(w, req)
		if err == nil {
			ret, err = handlerFunc(w, req, params)
		} else {
			ret = map[string]interface{}{"error": err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		var status int
		if err == nil {
			status = http.StatusOK
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)

		// debug request body content
		//log.Trace("req body: %+v", params)

		// access log
		log.Debug("%s \"%s %s %s\" %d %s",
			req.RemoteAddr,
			req.Method,
			req.RequestURI,
			req.Proto,
			status,
			time.Since(t1))
		if status != http.StatusOK {
			log.Error("HTTP: %v", err)
		}

		if ret != nil {
			// pretty write json result
			pretty, err := json.MarshalIndent(ret, "", "    ")
			if err != nil {
				log.Error(err)
				return
			}
			w.Write(pretty)
			w.Write([]byte("\n"))
		}
	}

	// path can't be duplicated
	isDup := false
	for _, p := range this.api.httpPaths {
		if p == path {
			log.Error("REST[%s] already registered", path)
			isDup = true
			break
		}
	}

	if !isDup {
		this.api.httpPaths = append(this.api.httpPaths, path)
	}

	return this.api.httpRouter.HandleFunc(path, wrappedFunc)
}

func (this *HttpJsonServer) UnregisterAllHttpApi() {
	this.api.httpPaths = this.api.httpPaths[:0]
}

func (this *httpRestApi) decodeHttpParams(w http.ResponseWriter,
	req *http.Request) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return params, nil
}
