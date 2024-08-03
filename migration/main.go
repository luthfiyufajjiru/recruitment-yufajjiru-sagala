package main

import (
	_ "embed"
	"sagala-todo/dependency"
	"sagala-todo/pkg/constants"
	customlog "sagala-todo/pkg/custom-log"
)

//go:embed table_init.sql
var script string

var logger = customlog.Logger

func main() {
	cfg := dependency.InitConfiguration()
	sqlInstances := dependency.InitMigration(cfg)

	logger.Info(script)

	_, err := sqlInstances[constants.ConnSqlDefault].Db.Exec(script)
	if err != nil {
		logger.Fatal(err)
	}

}
