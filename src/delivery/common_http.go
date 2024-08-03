package delivery

import (
	"fmt"
	"net/http"
	"sagala-todo/pkg/constants"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	UnknownHttpMethod = "unknown http method"
)

func HandleUnknownHttpMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
	fmt.Fprint(w, UnknownHttpMethod)
}

func HandleUnhandledError(w http.ResponseWriter, err error, logger *logrus.Entry) {
	errorId := uuid.NewString()
	logger.WithField(constants.LFKErrorId, errorId).Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
	fmt.Fprintf(w, "unhandled error, error id: %s", errorId)
}
