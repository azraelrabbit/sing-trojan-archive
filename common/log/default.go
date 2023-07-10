package log

import (
	"io"
	"time"

	"github.com/sagernet/sing/common/log/internal/default"
)

var _ = default_logger.Options(DefaultOptions{})

type DefaultOptions struct {
	Level          Level
	Writer         io.Writer
	BaseTime       time.Time
	PlatformWriter io.Writer
	DisableColor   bool
	Timestamp      bool
}

func NewDefault(options DefaultOptions) Logger {
	return default_logger.New(default_logger.Options(options))
}
