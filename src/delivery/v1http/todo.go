package v1http

import (
	"net/http"
	"sagala-todo/pkg/constants"
	"sagala-todo/src/delivery"
	"github.com/sirupsen/logrus"
)

func (h *V1Handler) RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
		case http.MethodGet:
		default:
			delivery.HandleUnknownHttpMethod(w)
			return
		}
	}
}

func (h *V1Handler) TaskDetail() http.HandlerFunc {
	const handlerName = "TaskDetail"
	logger := logger.WithFields(logrus.Fields{
		constants.LFKHandlerName: handlerName,
	})
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
		case http.MethodPut:
		case http.MethodDelete:
		default:
			delivery.HandleUnknownHttpMethod(w)
			return
		}
	}
}
