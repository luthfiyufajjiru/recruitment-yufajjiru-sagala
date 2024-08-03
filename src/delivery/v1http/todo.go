package v1http

import (
	"net/http"
	"sagala-todo/src/delivery"
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
