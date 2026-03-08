package handler

import (
	"fmt"
	"time"

	http "net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LogHandler struct {
}

func WithLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		h.ServeHTTP(res, req)

		duration := time.Since(start)

		requestID := uuid.New().String()

		logrus.WithFields(logrus.Fields{
			"url":          req.URL.Path,
			"method":       req.Method,
			"request_id":   requestID,
			"request_time": fmt.Sprintf("%vms", duration.Milliseconds()),
		}).Info("Request")

		logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"status":     res.Header().Get("Status"),
			"length":     res.Header().Get("Content-Length"),
		}).Info("Response")
	})
}
