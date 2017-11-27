package utsusemi

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ResponseWriterWithCode struct {
	http.ResponseWriter
	Code int
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
		origInfo := fmt.Sprintf("%s %s", request.Method, request.RequestURI)
		server.Logger.Print(origInfo)

		var proxy *httputil.ReverseProxy
		var url *url.URL
		backendLen := len(server.Backends)

	ScanLoop:
		for i, b := range server.Backends {
			proxy = b.Proxy
			url = b.URL

			if i == backendLen-1 {
				break ScanLoop
			}

			resp, err := http.Head(singleJoiningSlash(url.String(), request.RequestURI))

			if err != nil {
				continue
			}

			server.Logger.Printf("%s -> %s %s %d", origInfo, "HEAD", url.Host, resp.StatusCode)

			for _, ok := range b.Ok {
				if resp.StatusCode == ok {
					break ScanLoop
				}
			}
		}

		request.Host = url.Host
		proxy.ServeHTTP(writer, request)

		server.Logger.Printf("%s -> %s %s", origInfo, request.Method, request.Host)
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)

	return
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
