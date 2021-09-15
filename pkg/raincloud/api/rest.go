package api

import "net/http"

type RestServer struct {
	client http.Client
}

// NewRestServer initializes a new rest API
func NewRestServer() (*RestServer, error) {
	panic("implement me!")
}

// Run launches the main loop
func (rest *RestServer) Run() {
	panic("implement me!")
}

// Stop kills the server
func (rest *RestServer) Stop() {
	panic("implement me!")
}
