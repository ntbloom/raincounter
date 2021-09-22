package api

import (
	"net/http"

	"github.com/ntbloom/raincounter/pkg/config/configkey"

	"github.com/sirupsen/logrus"
)

type RestServer struct {
	server *http.Server
	mux    *http.ServeMux
	state  chan int
}

const (
	teapot   = "/v1.0/teapot"
	hello    = "/v1.0/hello"
	lastRain = "/v1.0/lastRain"
	rain     = "/v1.0/rain"
)

// NewRestServer initializes a new rest API
func NewRestServer() (*RestServer, error) {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	state := make(chan int, 1)
	return &RestServer{
		server: &server,
		mux:    mux,
		state:  state,
	}, nil
}

// Run launches the main loop
func (rest *RestServer) Run() {
	handler := newRestHandler()
	defer handler.close()

	// divert all endpoints to the handler
	rest.mux.Handle("/", handler)

	go logrus.Fatalf("problem with ListenAndServe: %s", rest.server.ListenAndServe())
	for {
		state := <-rest.state
		switch state {
		case configkey.Kill:
			err := rest.server.Close()
			if err != nil {
				logrus.Fatalf("problem closing server: %s", err)
			}
		default:
			logrus.Errorf("unexpected message on rest.state channel: %d", state)
		}
	}
}

// Stop kills the server
func (rest *RestServer) Stop() {
	logrus.Info("killing the rest API server")
	rest.state <- configkey.Kill
}
