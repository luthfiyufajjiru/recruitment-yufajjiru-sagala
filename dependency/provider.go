package dependency

import (
	"fmt"
	"os"
	"os/signal"
	"sagala-todo/pkg/adapters"
	"sagala-todo/pkg/common"
	"syscall"
)

func provideConfiguration() (c adapters.Config) {
	c = make(adapters.Config)
	c.Load(fmt.Sprintf("%s/.env", common.RootDirectory()))
	return
}

func provideSql(adCfg adapters.Config, cfgs []adapters.SqlConfig, dsns []string) (sql map[string]*adapters.Sql) {
	sql = make(map[string]*adapters.Sql)
	for i, cfg := range cfgs {
		cfg.Dsn = adCfg[dsns[i]]
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
