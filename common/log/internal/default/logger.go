package default_logger

import (
	"context"
	"io"
	"os"
	"runtime"
	"time"

	F "github.com/sagernet/sing/common/format"
	log "github.com/sagernet/sing/common/log"
)

type Options struct {
	Level          log.Level
	Writer         io.Writer
	BaseTime       time.Time
	PlatformWriter io.Writer
	DisableColor   bool
	Timestamp      bool
}

type abstractLogger struct {
	formatter         Formatter
	platformFormatter Formatter
	writer            io.Writer
	platformWriter    io.Writer
}

type Logger struct {
	*abstractLogger
	prefix string
	level  log.Level
}

func New(options Options) log.Logger {
	aLogger := &abstractLogger{
		formatter: Formatter{
			BaseTime:         options.BaseTime,
			DisableColors:    options.DisableColor,
			DisableTimestamp: !options.Timestamp,
			FullTimestamp:    true,
		},
		platformFormatter: Formatter{
			BaseTime:         options.BaseTime,
			DisableColors:    runtime.GOOS == "darwin" || runtime.GOOS == "ios",
			DisableLineBreak: true,
		},
		writer: options.Writer,
	}
	return &Logger{
		abstractLogger: aLogger,
		level:          options.Level,
	}
}

func (l *Logger) WithPrefix(prefix string) log.Logger {
	return &Logger{
		abstractLogger: l.abstractLogger,
		prefix:         prefix,
		level:          l.level,
	}
}

func (l *Logger) WithLevel(level log.Level) log.Logger {
	return &Logger{
		abstractLogger: l.abstractLogger,
		prefix:         l.prefix,
		level:          level,
	}
}

func (l *Logger) Log(ctx context.Context, level log.Level, args []any) {
	if level > l.level {
		return
	}
	nowTime := time.Now()
	if l.writer != nil {
		_, _ = l.writer.Write([]byte(l.formatter.Format(ctx, level, l.prefix, F.ToString(args...), nowTime)))
	}
	if l.platformWriter != nil {
		_, _ = l.platformWriter.Write([]byte(l.platformFormatter.Format(ctx, level, l.prefix, F.ToString(args...), nowTime)))
	}
	if level == log.LevelPanic {
		panic(F.ToString(args...))
	}
	if level == log.LevelFatal {
		os.Exit(1)
	}
}
