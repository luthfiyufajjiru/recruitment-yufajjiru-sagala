package dependency

import (
	"os"
	"os/signal"
	"sagala-todo/pkg/adapters"
	"sagala-todo/pkg/common"
	"syscall"
)

func provideConfiguration() (c adapters.Config) {
	c = make(adapters.Config)
	c.Load(common.RootDirectory())
	return
}

func provideSql(cfgs []adapters.SqlConfig) (sql map[string]*adapters.Sql) {
	sql = make(map[string]*adapters.Sql)
	for _, cfg := range cfgs {
		s := new(adapters.Sql)
		s.Init(&cfg)
		sql[cfg.RegistryName] = s
	}
	return
}

func provideOsSignal() chan os.Signal {
	terminalHandler := make(chan os.Signal, 1)
	signal.Notify(terminalHandler,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	return terminalHandler
}
