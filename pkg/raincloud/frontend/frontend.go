package frontend

import (
	"net/http"
	"text/template"

	"github.com/ntbloom/raincounter/pkg/raincloud/frontend/fetch"

	"github.com/spf13/viper"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/sirupsen/logrus"
)

type HTMLServer struct {
	server  *http.Server
	mux     *http.ServeMux
	fetcher fetch.Fetcher
	state   chan int
}

func NewHTMLServer() (*HTMLServer, error) {
	webAddress := viper.GetString(configkey.WebServerAddress)
	logrus.Debugf("webAddress=%s", webAddress)
	mux := http.NewServeMux()
	server := http.Server{
		Addr:              webAddress,
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

	return &HTMLServer{
		server:  &server,
		mux:     mux,
		fetcher: fetch.NewDataFetcher(),
		state:   state,
	}, nil
}

// Run serves static html pages
func (h *HTMLServer) Run() {
	logrus.Infof("starting the server on %s", h.server.Addr)
	var entrypoint = viper.GetString(configkey.WebEntrypoint)
	logrus.Debugf("entrypoint=%s", entrypoint)
	tmpl := template.Must(template.ParseFiles(entrypoint))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := h.fetcher.Fetch()
		err := tmpl.Execute(w, data)
		if err != nil {
			logrus.Errorf("error fetching data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	go func() {
		err := http.ListenAndServe(h.server.Addr, nil)
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	for {
		state := <-h.state
		switch state {
		case configkey.Kill:
			err := h.server.Close()
			if err != nil {
				logrus.Fatalf("problem closing server: %s", err)
			}
		default:
			logrus.Errorf("unexpected message on rest.state channel: %d", state)
		}
	}
}

// Stop kills the html server
func (h *HTMLServer) Stop() {
	logrus.Info("killing the rest API server")
	h.state <- configkey.Kill
}
