package cmd

import (
	"net/http"
)

func InitHttpServer() {
	mux := http.NewServeMux()

	// v1Todo := dependency.InitTodoV1HttpHandler()

	// mux.Handle("/v1", v1Todo.RootHandler())
	// mux.Handle("/v1/{id}", v1Todo.TaskDetail())

	http.ListenAndServe("", mux)
}
