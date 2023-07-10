package log

import "context"

type Logger interface {
	WithPrefix(prefix string) Logger
	WithLevel(level Level) Logger

	Trace(args ...any)
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Panic(args ...any)

	TraceContext(ctx context.Context, args ...any)
	DebugContext(ctx context.Context, args ...any)
	InfoContext(ctx context.Context, args ...any)
	WarnContext(ctx context.Context, args ...any)
	ErrorContext(ctx context.Context, args ...any)
	FatalContext(ctx context.Context, args ...any)
	PanicContext(ctx context.Context, args ...any)
}
