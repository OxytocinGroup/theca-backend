package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(ctx context.Context, message string, fields map[string]any)
	Info(ctx context.Context, message string, fields map[string]any)
	Warn(ctx context.Context, message string, fields map[string]any)
	Error(ctx context.Context, message string, fields map[string]any)
}

type LogrusLogger struct {
	log *logrus.Logger
}

func NewLogrusLogger(level string) *LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	infoFile := &lumberjack.Logger{
		Filename:   "logs/info.log",
		MaxSize:    10, // Размер в MB
		MaxBackups: 3,
		MaxAge:     30, // Возраст в днях
		Compress:   true,
	}
	warnFile := &lumberjack.Logger{
		Filename:   "logs/warn.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}
	errorFile := &lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}
	logger.SetOutput(os.Stdout)
	logger.AddHook(NewLevelHook(infoFile, logrus.InfoLevel))
	logger.AddHook(NewLevelHook(warnFile, logrus.WarnLevel))
	logger.AddHook(NewLevelHook(errorFile, logrus.ErrorLevel))

	return &LogrusLogger{log: logger}
}

type LevelHook struct {
	writer    *lumberjack.Logger
	logLevels []logrus.Level
}

func NewLevelHook(writer *lumberjack.Logger, levels ...logrus.Level) *LevelHook {
	return &LevelHook{
		writer:    writer,
		logLevels: levels,
	}
}

func (h *LevelHook) Levels() []logrus.Level {
	return h.logLevels
}

func (h *LevelHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = h.writer.Write([]byte(line))
	return err
}

func (l *LogrusLogger) Debug(ctx context.Context, message string, fields map[string]any) {
	l.log.WithFields(logrus.Fields(fields)).Debug(message)
}

func (l *LogrusLogger) Info(ctx context.Context, message string, fields map[string]any) {
	l.log.WithFields(logrus.Fields(fields)).Info(message)
}

func (l *LogrusLogger) Warn(ctx context.Context, message string, fields map[string]any) {
	l.log.WithFields(logrus.Fields(fields)).Warn(message)
}

func (l *LogrusLogger) Error(ctx context.Context, message string, fields map[string]any) {
	l.log.WithFields(logrus.Fields(fields)).Error(message)
}
