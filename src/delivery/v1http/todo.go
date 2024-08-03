package v1http

import "net/http"

func (h *V1Handler) RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *V1Handler) TaskDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
