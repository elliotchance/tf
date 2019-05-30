package tf

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type HTTPServer struct {
	Server   *http.Server
	Port     int
	Shutdown func()
	Mux      *http.ServeMux
}

func (server *HTTPServer) Endpoint() string {
	return fmt.Sprintf("http://localhost:%d", server.Port)
}

func (server *HTTPServer) AddHandler(path string, handler http.HandlerFunc) *HTTPServer {
	// Path cannot be blank, otherwise Mux.HandleFunc will panic.
	if path == "" {
		path = "/"
	}

	server.Mux.HandleFunc(path, handler)

	return server
}

func (server *HTTPServer) AddHandlers(handlers map[string]http.HandlerFunc) *HTTPServer {
	for path, handler := range handlers {
		server.AddHandler(path, handler)
	}

	return server
}

func StartHTTPServer(port int) *HTTPServer {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	server := &HTTPServer{
		Mux:    mux,
		Port:   port,
		Server: srv,
		Shutdown: func() {
			// In some cases the shutdown will panic. We don't care about
			// graceful shutdown under test.
			defer func() {
				recover()
			}()

			_ = srv.Shutdown(nil)
		},
	}

	// ListenAndServe() is not safe under test because it's possible the test
	// will make a request before the listener is setup. So split them into
	// their separate steps.
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		panic(err)
	}

	listeningOn := strings.Split(listener.Addr().String(), ":")
	server.Port, _ = strconv.Atoi(listeningOn[len(listeningOn)-1])

	go func() {
		// This will always return the error "http: Server closed" because the
		// test explicitly closes it.
		err := srv.Serve(listener)

		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()

	return server
}

func HTTPEmptyResponse(statusCode int) func(http.ResponseWriter, *http.Request) {
	return HTTPStringResponse(statusCode, "")
}

func HTTPStringResponse(statusCode int, body string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, err := w.Write([]byte(body))
		if err != nil {
			panic(err)
		}
	}
}

func HTTPJSONResponse(statusCode int, body interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		data, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		_, err = w.Write(data)
		if err != nil {
			panic(err)
		}
	}
}
