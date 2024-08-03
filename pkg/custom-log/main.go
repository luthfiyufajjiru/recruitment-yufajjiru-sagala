package customlog

import (
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetReportCaller(true)
}

func SetLevel(expectedLevel string) {
	var level logrus.Level
	switch expectedLevel {
	case "info":
		level = logrus.InfoLevel
	case "panic":
		level = logrus.PanicLevel
	case "debug":
		level = logrus.DebugLevel
	case "error":
		level = logrus.ErrorLevel
	default:
		panic("not known log level")
	}

	Logger.SetLevel(level)
}

func SetLoki(appName, env, uri string) {
	opts := lokirus.NewLokiHookOptions().
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithFormatter(&logrus.JSONFormatter{}).
		WithStaticLabels(lokirus.Labels{
			"app":         appName,
			"environment": env,
		})

	hook := lokirus.NewLokiHookWithOpts(
		uri,
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	)

	Logger.AddHook(hook)
}
