package utsusemi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

type ResponseWriterWithCode struct {
	http.ResponseWriter
	Code int
}

func NewResponseWriterWithCode(w http.ResponseWriter) *ResponseWriterWithCode {
	return &ResponseWriterWithCode{w, http.StatusOK}
}

func (w *ResponseWriterWithCode) WriteHeader(code int) {
	w.Code = code
	w.ResponseWriter.WriteHeader(code)
}

type Backend struct {
	Proxy *httputil.ReverseProxy
	URL   *url.URL
	Ok    []int
}

type Server struct {
	Port     int
	Backends []Backend
	Logger   *log.Logger
}

func NewServer(config *Config, logger *log.Logger) (server *Server, err error) {
	backends := make([]Backend, len(config.Backend))

	for i, b := range config.Backend {
		var target *url.URL
		target, err = url.Parse(b.Target)

		if err != nil {
			return
		}

		backends[i].Proxy = httputil.NewSingleHostReverseProxy(target)
		backends[i].URL = target
		backends[i].Ok = b.Ok
	}

	server = &Server{
		Port:     config.Port,
		Backends: backends,
		Logger:   logger,
	}

	return
}

func (server *Server) Run() (err error) {
	backendLen := len(server.Backends)

	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "pong")
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		for i := 0; i < backendLen; i++ {
			backend := server.Backends[i]
			request.Host = backend.URL.Host

			if i < backendLen-1 {
				preWriter := httptest.NewRecorder()
				backend.Proxy.ServeHTTP(preWriter, request)
				server.logRequest(request, preWriter.Code)

				if isResponseOk(preWriter, backend.Ok) {
					io.Copy(writer, preWriter.Body)
					break
				}
			} else {
				writerWithCode := NewResponseWriterWithCode(writer)
				backend.Proxy.ServeHTTP(writerWithCode, request)
				server.logRequest(request, writerWithCode.Code)
			}
		}
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)

	return
}

func isResponseOk(res *httptest.ResponseRecorder, okCodes []int) bool {
	for _, ok := range okCodes {
		if res.Code == ok {
			return true
		}
	}

	return false
}

func (server *Server) logRequest(request *http.Request, code int) {
	server.Logger.Printf("%s %s -> %s %s %d", request.Method, request.RequestURI, request.Method, request.Host, code)
}
