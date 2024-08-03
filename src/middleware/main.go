package middleware

import (
	customlog "sagala-todo/pkg/custom-log"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger = customlog.Logger
