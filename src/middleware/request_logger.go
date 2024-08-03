package middleware

import (
	"net/http"
	"sagala-todo/pkg/constants"

	"github.com/sirupsen/logrus"
)

func RequestLogger(next http.Handler) http.HandlerFunc {
	logger := logger.WithFields(logrus.Fields{
		constants.LFKHandlerName: "request_logger_middleware",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}
