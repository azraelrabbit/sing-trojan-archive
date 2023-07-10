package default_logger

import (
	"context"
	"strconv"
	"strings"
	"time"

	F "github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/common/log"
)

type Formatter struct {
	BaseTime         time.Time
	DisableColors    bool
	DisableTimestamp bool
	FullTimestamp    bool
	DisableLineBreak bool
}

func (f Formatter) Format(ctx context.Context, level log.Level, tag string, message string, timestamp time.Time) string {
	const timeFormat = "-0700 2006-01-02 15:04:05"
	levelString := strings.ToUpper(log.FormatLevel(level))
	if !f.DisableColors {
		var levelColor uint
		switch level {
		case log.LevelDebug, log.LevelTrace:
			levelColor = uint(WhiteFg)
		case log.LevelInfo:
			levelColor = uint(CyanFg)
		case log.LevelWarn:
			levelColor = uint(YellowFg)
		case log.LevelError, log.LevelFatal, log.LevelPanic:
			levelColor = uint(RedFg)
		}
		levelString = Colorize(levelColor, levelString)
	}
	if tag != "" {
		message = tag + ": " + message
	}
	var id log.ContextID
	var hasId bool
	if ctx != nil {
		id, hasId = log.IDFromContext(ctx)
	}
	if hasId {
		activeDuration := formatDuration(time.Since(id.CreatedAt))
		if !f.DisableColors {
			var color uint
			color = uint(uint8(id.ID))
			color %= 215
			row := color / 36
			column := color % 36

			var r, g, b float32
			r = float32(row * 51)
			g = float32(column / 6 * 51)
			b = float32((column % 6) * 51)
			luma := 0.2126*r + 0.7152*g + 0.0722*b
			if luma < 60 {
				row = 5 - row
				column = 35 - column
				color = row*36 + column
			}
			color += 16
			color = color << 16
			color |= 1 << 14

			message = F.ToString("[", Colorize(color, id.ID), " ", activeDuration, "] ", message)
		} else {
			message = F.ToString("[", id.ID, " ", activeDuration, "] ", message)
		}
	}
	switch {
	case f.DisableTimestamp:
		message = levelString + " " + message
	case f.FullTimestamp:
		message = timestamp.Format(timeFormat) + " " + levelString + " " + message
	default:
		message = levelString + "[" + xd(int(timestamp.Sub(f.BaseTime)/time.Second), 4) + "] " + message
	}
	if f.DisableLineBreak {
		if message[len(message)-1] == '\n' {
			message = message[:len(message)-1]
		}
	} else {
		if message[len(message)-1] != '\n' {
			message += "\n"
		}
	}
	return message
}

func xd(value int, x int) string {
	message := strconv.Itoa(value)
	for len(message) < x {
		message = "0" + message
	}
	return message
}

func formatDuration(duration time.Duration) string {
	if duration < time.Second {
		return F.ToString(duration.Milliseconds(), "ms")
	} else if duration < time.Minute {
		return F.ToString(int64(duration.Seconds()), ".", int64(duration.Seconds()*100)%100, "s")
	} else {
		return F.ToString(int64(duration.Minutes()), "m", int64(duration.Seconds())%60, "s")
	}
}

func Colorize(color uint, message any) string {
	return "\033[" + Color(color).Nos(false) + "m" + F.ToString(message) + "0m"
}
