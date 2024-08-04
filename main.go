package main

import (
	"context"
	"net/http"
	"os"
	"sagala-todo/cmd"
	"sagala-todo/dependency"
	customlog "sagala-todo/pkg/custom-log"
	"sync"

	"github.com/sirupsen/logrus"
)

func signalDisrupt(cancel context.CancelFunc, ch chan os.Signal, logger *logrus.Logger) {
	defer cancel()
	data := <-ch
	logger.Infof("system call: %+v", data)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger := customlog.Logger

	osCh := dependency.InitOsSignalChannel()
	defer close(osCh)

	var wg sync.WaitGroup

	cfg := dependency.InitConfiguration()
	var httpServer http.Server

	go signalDisrupt(cancel, osCh, logger)
	go cmd.InitHttpServer(ctx, cfg, &wg, &httpServer)

	<-ctx.Done() // waiting for signal from os signal
	wg.Wait()    // wait for application to shut down
	logger.Info("app: bye")
}
