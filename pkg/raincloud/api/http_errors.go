package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// generic bad request handler
func (handler restHandler) badRequest(w http.ResponseWriter, err error) {
	logrus.Error(err)
	w.WriteHeader(http.StatusBadRequest)
}

// generic internal service error handler
func (handler restHandler) internalServiceError(w http.ResponseWriter, err error) {
	logrus.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
}
