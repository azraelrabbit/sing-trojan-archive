package log

import (
	"context"
	"os"
	"time"
)

var defaultLogger Logger

func init() {
	defaultLogger = NewDefault(DefaultOptions{
		Level:    LevelDebug,
		Writer:   os.Stderr,
		BaseTime: time.Now(),
	})
}

func Default() Logger {
	return defaultLogger
}

func SetDefault(logger Logger) {
	defaultLogger = logger
}

func Trace(args ...any) {
	defaultLogger.Trace(args...)
}

func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

func Info(args ...any) {
	defaultLogger.Info(args...)
}

func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

func Error(args ...any) {
	defaultLogger.Error(args...)
}

func Fatal(args ...any) {
	defaultLogger.Fatal(args...)
}

func Panic(args ...any) {
	defaultLogger.Panic(args...)
}

func TraceContext(ctx context.Context, args ...any) {
	defaultLogger.TraceContext(ctx, args...)
}

func DebugContext(ctx context.Context, args ...any) {
	defaultLogger.DebugContext(ctx, args...)
}

func InfoContext(ctx context.Context, args ...any) {
	defaultLogger.InfoContext(ctx, args...)
}

func WarnContext(ctx context.Context, args ...any) {
	defaultLogger.WarnContext(ctx, args...)
}

func ErrorContext(ctx context.Context, args ...any) {
	defaultLogger.ErrorContext(ctx, args...)
}

func FatalContext(ctx context.Context, args ...any) {
	defaultLogger.FatalContext(ctx, args...)
}

func PanicContext(ctx context.Context, args ...any) {
	defaultLogger.PanicContext(ctx, args...)
}
