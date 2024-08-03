package adapters

import (
	customlog "sagala-todo/pkg/custom-log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger = customlog.Logger

type (
	SqlConfig struct {
		RegistryName, DriverName, Dsn string
		MaxIdleTime, MaxLifeTime      time.Duration
		MaxIdleConns, MaxOpenConns    int
	}

	Sql struct {
		Db       *sqlx.DB
		MockCtrl sqlmock.Sqlmock
	}
)

func (s *Sql) Init(cfg *SqlConfig) {
	sqlDB, err := sqlx.Open(cfg.DriverName, cfg.Dsn)
	if err != nil {
		logger.Panic("error occurred while connecting with the database")
	}

	if cfg.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Minute * cfg.MaxIdleTime)
	}

	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.MaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime * time.Hour)
	}

	if err = sqlDB.Ping(); err != nil {
		logger.WithError(err).Panic("error occurred while connecting with the database")
	}

	s.Db = sqlDB
}
