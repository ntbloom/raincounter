package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

/* Simple functions to respond to requests with 400/500-level response codes */

func (handler restHandler) badRequest(w http.ResponseWriter, err error) {
	logrus.Error(err)
	w.WriteHeader(http.StatusBadRequest)
}

func (handler restHandler) internalServiceError(w http.ResponseWriter, err error) {
	logrus.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (handler restHandler) unsupportedMedia(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnsupportedMediaType)
}

func (handler restHandler) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func (handler restHandler) methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
