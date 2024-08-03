package cmd

import (
	"context"
	"fmt"
	"net/http"
	"sagala-todo/dependency"
	"sagala-todo/pkg/adapters"
	"sagala-todo/pkg/constants"
	customlog "sagala-todo/pkg/custom-log"
	"sagala-todo/src/middleware"
	"sync"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger = customlog.Logger

func InitHttpServer(ctx context.Context, cfg adapters.Config, wg *sync.WaitGroup, server *http.Server) {
	wg.Add(1)

	addr := fmt.Sprintf(cfg[constants.ADDR])
	mux := http.NewServeMux()

	*server = http.Server{
		Addr:    addr,
		Handler: mux,
	}

	v1Todo := dependency.InitTodoV1HttpHandler(cfg)

	mux.Handle("/v1", middleware.RequestLogger(v1Todo.RootHandler()))
	mux.Handle("/v1/{id}", middleware.RequestLogger(v1Todo.RootHandler()))

	go func(wg *sync.WaitGroup) {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			logger.Println("error during shutdown:", err)
			return
		}
		logger.Println("http server gracefully shut down")
		wg.Done()
	}(wg)

	logger.Printf("http Server started on %s", addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("ListenAndServe:", err)
	}
}
