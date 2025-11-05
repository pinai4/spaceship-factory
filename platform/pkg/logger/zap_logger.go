package logger

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	origZapLogger *zap.Logger
}

func NewZapLogger(levelStr string, asJSON bool) *zapLogger {
	dynamicLevel := zap.NewAtomicLevelAt(parseLevel(levelStr))

	encoderCfg := buildProductionEncoderConfig()

	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		dynamicLevel,
	)

	origZapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return &zapLogger{
		origZapLogger: origZapLogger,
	}
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",                 // time
		LevelKey:       "level",                     // log level
		NameKey:        "logger",                    // logger name, if used
		CallerKey:      "caller",                    // where the log was called from
		MessageKey:     "message",                   // message text
		StacktraceKey:  "stacktrace",                // stack trace for errors
		LineEnding:     zapcore.DefaultLineEnding,   // line break
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // readable ISO 8601 format
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // short caller
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	allFields := append(fieldsFromContext(ctx), convertFields(fields)...)
	l.origZapLogger.Debug(msg, allFields...)
}

func (l *zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	allFields := append(fieldsFromContext(ctx), convertFields(fields)...)
	l.origZapLogger.Info(msg, allFields...)
}

func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	allFields := append(fieldsFromContext(ctx), convertFields(fields)...)
	l.origZapLogger.Warn(msg, allFields...)
}

func (l *zapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	allFields := append(fieldsFromContext(ctx), convertFields(fields)...)
	l.origZapLogger.Error(msg, allFields...)
}

func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	allFields := append(fieldsFromContext(ctx), convertFields(fields)...)
	l.origZapLogger.Fatal(msg, allFields...)
}

// parseLevel convert string level to zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsFromContext extract enrich-fields from context
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}

// internal conversion
func convertFields(fields []Field) []zap.Field {
	zf := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		switch v := f.Value.(type) {
		case string:
			zf = append(zf, zap.String(f.Key, v))
		case int:
			zf = append(zf, zap.Int(f.Key, v))
		case int32:
			zf = append(zf, zap.Int32(f.Key, v))
		case int64:
			zf = append(zf, zap.Int64(f.Key, v))
		case error:
			zf = append(zf, zap.Error(v))
		default:
			zf = append(zf, zap.Any(f.Key, v))
		}
	}
	return zf
}
