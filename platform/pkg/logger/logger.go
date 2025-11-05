package logger

import (
	"context"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}

func String(key, val string) Field {
	return Field{Key: key, Value: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Value: val}
}

func Int32(key string, val int32) Field {
	return Field{Key: key, Value: val}
}

func Int64(key string, val int64) Field {
	return Field{Key: key, Value: val}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}
