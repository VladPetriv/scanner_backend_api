package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

var (
	logger Logger    // nolint
	once   sync.Once // nolint
)

func Get(logLevel string) *Logger {
	once.Do(func() {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoder := zapcore.NewJSONEncoder(config)
		consoleEncoder := zapcore.NewConsoleEncoder(config)

		err := os.MkdirAll("logs", 0o755) // nolint
		if err != nil {
			panic(err)
		}

		logFile, err := os.OpenFile("logs/logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644) // nolint
		if err != nil {
			panic(err)
		}

		writer := zapcore.AddSync(logFile)

		logLevel, err := zapcore.ParseLevel(logLevel)
		if err != nil {
			panic(err)
		}

		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, logLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
		)

		newLoggger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

		logger = Logger{newLoggger}
	})

	return &logger
}
