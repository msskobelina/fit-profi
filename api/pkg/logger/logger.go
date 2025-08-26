package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Named(name string) Interface
	With(args ...interface{}) Interface
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

type Logger struct {
	logger *zap.SugaredLogger
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	var l zapcore.Level
	switch strings.ToLower(level) {
	case "error":
		l = zapcore.ErrorLevel
	case "warn":
		l = zapcore.WarnLevel
	case "info":
		l = zapcore.InfoLevel
	case "debug":
		l = zapcore.DebugLevel
	default:
		l = zapcore.InfoLevel
	}

	config := zap.Config{
		Development:      false,
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(l),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeDuration: zapcore.SecondsDurationEncoder,
			LevelKey:       "severity",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			TimeKey:        "timestamp",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			NameKey:        "name",
			EncodeName:     zapcore.FullNameEncoder,
			MessageKey:     "message",
			StacktraceKey:  "",
			LineEnding:     "\n",
		},
	}

	logger, _ := config.Build()

	return &Logger{
		logger: logger.Sugar(),
	}
}

func (l *Logger) Named(name string) Interface {
	return &Logger{
		logger: l.logger.Named(name),
	}
}

func (l *Logger) With(args ...interface{}) Interface {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debugw(message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Infow(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warnw(message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Errorw(message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Fatalw(message, args...)
	os.Exit(1)
}
