package api

import "net/http"

// handleOptions handles OPTIONS requests
func (handler restHandler) handleOptions(w http.ResponseWriter) {
	for k, v := range map[string]string{
		"Allow":                         "OPTIONS, GET",
		"Access-Control-Request-Method": "GET",
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Allow-Headers":  "Content-Type",
	} {
		w.Header().Set(k, v)
	}
	w.WriteHeader(http.StatusOK)
}
