package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger = *logrus.Entry

func NewLogger() Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "timestamp",
		},

		TimestampFormat: time.RFC3339Nano,
	})

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}

	log.SetLevel(level)

	log.SetOutput(os.Stdout)

	logger := log.WithFields(logrus.Fields{
		"service": os.Getenv("SERVICE_NAME"),
		"version": os.Getenv("APP_VERSION"),
	})

	return logger
}

/* example you can see this answer in the log

{
	"message": "request suuccess",
	"severity": "info",
	"timestamp": "2023-05-01T12:00:00.000Z",
	"service": "CashFlow",
	"version": "1.0.0"
  }
*/
