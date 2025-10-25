package logger

import (
	"context"
)

type NoopLogger struct{}

func (l *NoopLogger) Info(ctx context.Context, msg string, fields ...Field)  {}
func (l *NoopLogger) Error(ctx context.Context, msg string, fields ...Field) {}
