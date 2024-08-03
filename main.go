package main

import (
	"os"
	"sagala-todo/cmd"
	"sagala-todo/dependency"
	customlog "sagala-todo/pkg/custom-log"

	"github.com/sirupsen/logrus"
)

func signalDisrupt(ch chan os.Signal, logger *logrus.Logger) {
	data := <-ch
	logger.Errorf("system call: %+v\n", data)
}

func main() {
	logger := customlog.Logger

	osCh := dependency.InitOsSignalChannel()

	go signalDisrupt(osCh, logger)
	go cmd.InitHttpServer()
	go cmd.InitAutomationTest()
}
