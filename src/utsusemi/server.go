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
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		backendLen := len(server.Backends)

		for i, b := range server.Backends {
			proxy := b.Proxy
			url := b.URL

			if i < backendLen-1 {
				preWriter := httptest.NewRecorder()
				preRequest := httptest.NewRequest(request.Method, request.RequestURI, request.Body)
				preRequest.Host = url.Host
				proxy.ServeHTTP(preWriter, preRequest)
				server.logRequest(request, preRequest.Method, preRequest.Host, preWriter.Code)

				if isResponseOk(preWriter, b.Ok) {
					io.Copy(writer, preWriter.Body)
					break
				}
			} else {
				request.Host = url.Host
				writerWithCode := NewResponseWriterWithCode(writer)
				proxy.ServeHTTP(writerWithCode, request)
				server.logRequest(request, request.Method, request.Host, writerWithCode.Code)
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

func (server *Server) logRequest(origRequest *http.Request, method string, host string, code int) {
	server.Logger.Printf("%s %s -> %s %s %d", origRequest.Method, origRequest.RequestURI, method, host, code)
}
