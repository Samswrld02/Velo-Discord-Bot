package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() (*logrus.Logger, *os.File, error) {
	file, err := os.OpenFile("./logger/logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, nil, err
	}

	logger := logrus.New()

	logger.SetOutput(io.MultiWriter(os.Stdout, file))
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return logger, file, nil
}