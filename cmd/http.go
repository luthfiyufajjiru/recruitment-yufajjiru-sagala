package cmd

import (
	"fmt"
	"net/http"
	"sagala-todo/dependency"
	"sagala-todo/pkg/adapters"
	"sagala-todo/pkg/constants"
)

func InitHttpServer(cfg adapters.Config) {
	mux := http.NewServeMux()

	v1Todo := dependency.InitTodoV1HttpHandler(cfg)

	mux.Handle("/v1", v1Todo.RootHandler())
	mux.Handle("/v1/{id}", v1Todo.TaskDetail())

	http.ListenAndServe(fmt.Sprintf(":%s", cfg[constants.PORT]), mux)
}
