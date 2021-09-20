package api

import (
	"encoding/json"
	"net/http"

	"github.com/ntbloom/raincounter/pkg/config/configkey"

	"github.com/sirupsen/logrus"
)

const (
	contentType = "Content-Type"
	appJson     = "application/json"
)

type RestServer struct {
	server *http.Server
	mux    *http.ServeMux
	state  chan int
}

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
	rest.mux.HandleFunc("/v1.0/teapot", handleTeapot)
	rest.mux.HandleFunc("/v1.0/hello", handleHello)

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

// return teapot messages as bellweather for general server and for bootstrapping
// may be able to delete this later as the API is developed
func handleTeapot(w http.ResponseWriter, res *http.Request) {
	var payload []byte
	var err error

	encoding := res.Header.Get(contentType)
	logrus.Debugf("recevied request with `%s` encoding", encoding)
	if payload, err = json.Marshal(map[string]string{"hello": "teapot"}); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusTeapot)
		if _, err = w.Write(payload); err != nil {
			logrus.Error(err)
		}
	}
}

// template for json payload messages
func handleHello(w http.ResponseWriter, res *http.Request) {
	var payload []byte
	var err error

	encoding := res.Header.Get(contentType)
	logrus.Debugf("recevied request with `%s` encoding", encoding)
	if encoding != appJson {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	} else {
		if payload, err = json.Marshal(map[string]string{"hello": "world"}); err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			if _, err = w.Write(payload); err != nil {
				logrus.Error(err)
			}
		}
	}
}
