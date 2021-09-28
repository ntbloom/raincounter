package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

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
