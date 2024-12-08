
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "os"
)

var Log *zap.Logger

// InitializeLogger sets up the zap logger with custom settings.
func InitializeLogger(logFilePath string) error {
    // Create a file to store logs
    file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        return err
    }

    // Create a core that writes logs to the file and console
    fileCore := zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        zapcore.AddSync(file),
        zapcore.InfoLevel,
    )

    consoleCore := zapcore.NewCore(
        zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
        zapcore.AddSync(os.Stdout),
        zapcore.DebugLevel,
    )

    // Combine file and console cores
    combinedCore := zapcore.NewTee(fileCore, consoleCore)

    // Create the logger instance
    Log = zap.New(combinedCore, zap.AddCaller())
    return nil
}

// Info logs informational messages
func Info(msg string, fields ...zap.Field) {
    Log.Info(msg, fields...)
}

// Warn logs warning messages
func Warn(msg string, fields ...zap.Field) {
    Log.Warn(msg, fields...)
}

// Error logs error messages
func Error(msg string, fields ...zap.Field) {
    Log.Error(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() {
    _ = Log.Sync()
}
