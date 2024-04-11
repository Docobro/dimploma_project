package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LegacyLogger struct {
	logger Logger
}

func NewLegacyLogger(logger Logger) *LegacyLogger {
	return &LegacyLogger{logger: logger}
}

func (l LegacyLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Info(format, toFields(v)...)
}

func (l LegacyLogger) Debug(args ...interface{}) {
	l.logger.Debug("", toFields(args)...)
}

func (l LegacyLogger) Info(args ...interface{}) {
	l.logger.Info("", toFields(args)...)
}

func (l LegacyLogger) Warn(args ...interface{}) {
	l.logger.Warn("", toFields(args)...)
}

func (l LegacyLogger) Error(args ...interface{}) {
	l.logger.Error("", toFields(args)...)
}

func (l LegacyLogger) Fatal(args ...interface{}) {
	l.logger.Fatal("", toFields(args)...)
}

// Convert arbitrary arguments to zapcore.Fields
func toFields(args []interface{}) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(args))
	for _, arg := range args {
		fields = append(fields, zap.Any("arg", arg))
	}
	return fields
}
