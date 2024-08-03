package delivery

import (
	"fmt"
	"net/http"
	"sagala-todo/pkg/constants"
)

const (
	UnknownHttpMethod = "unknown http method"
)

func HandleUnknownHttpMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set(constants.HeaderKeyContentType, constants.HeaderTextPlain)
	fmt.Fprint(w, UnknownHttpMethod)
}
