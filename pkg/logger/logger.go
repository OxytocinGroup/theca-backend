package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(ctx context.Context, message string, fields map[string]interface{})
	Info(ctx context.Context, message string, fields map[string]interface{})
	Warn(ctx context.Context, message string, fields map[string]interface{})
	Error(ctx context.Context, message string, fields map[string]interface{})
}

type LogrusLogger struct {
	log *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	return &LogrusLogger{log: logger}
}

func (l *LogrusLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.log.WithFields(logrus.Fields(fields)).Debug(message)
}

func (l *LogrusLogger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.log.WithFields(logrus.Fields(fields)).Info(message)
}
func (l *LogrusLogger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.log.WithFields(logrus.Fields(fields)).Warn(message)
}
func (l *LogrusLogger) Error(ctx context.Context, message string, fields map[string]interface{}) {
	l.log.WithFields(logrus.Fields(fields)).Error(message)
}
