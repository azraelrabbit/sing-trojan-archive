package default_logger

import (
	"context"

	"github.com/sagernet/sing/common/log"
)

func (l *Logger) Trace(args ...any) {
	l.Log(context.Background(), log.LevelTrace, args)
}

func (l *Logger) Debug(args ...any) {
	l.Log(context.Background(), log.LevelDebug, args)
}

func (l *Logger) Info(args ...any) {
	l.Log(context.Background(), log.LevelInfo, args)
}

func (l *Logger) Warn(args ...any) {
	l.Log(context.Background(), log.LevelWarn, args)
}

func (l *Logger) Error(args ...any) {
	l.Log(context.Background(), log.LevelError, args)
}

func (l *Logger) Fatal(args ...any) {
	l.Log(context.Background(), log.LevelFatal, args)
}

func (l *Logger) Panic(args ...any) {
	l.Log(context.Background(), log.LevelPanic, args)
}

func (l *Logger) TraceContext(ctx context.Context, args ...any) {
	l.Log(ctx, log.LevelTrace, args)
}

func (l *Logger) DebugContext(ctx context.Context, args ...any) {
	l.Log(ctx, log.LevelDebug, args)
}

func (l *Logger) InfoContext(ctx context.Context, args ...any) {
	l.Log(ctx, log.LevelInfo, args)
}

func (l *Logger) WarnContext(ctx context.Context, args ...any) {
	l.Log(ctx, log.LevelWarn, args)
}

func (l *Logger) ErrorContext(ctx context.Context, args ...any) {
	l.Log(ctx, log.LevelError, args)
}

func (l *Logger) FatalContext(ctx context.Context, args ...any) {
	l.Log(ctx, log.LevelFatal, args)
}

func (l *Logger) PanicContext(ctx context.Context, args ...any) {
	l.Log(ctx, log.LevelPanic, args)
}
