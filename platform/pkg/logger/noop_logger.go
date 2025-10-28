package logger

import (
	"context"
)

type NoopLogger struct{}

func (l *NoopLogger) Debug(ctx context.Context, msg string, fields ...Field) {}
func (l *NoopLogger) Info(ctx context.Context, msg string, fields ...Field)  {}
func (l *NoopLogger) Warn(ctx context.Context, msg string, fields ...Field)  {}
func (l *NoopLogger) Error(ctx context.Context, msg string, fields ...Field) {}
func (l *NoopLogger) Fatal(ctx context.Context, msg string, fields ...Field) {}
