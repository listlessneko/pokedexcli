package timber

import (
	"log/slog"
)

const (
	LevelSawdust  = slog.Level(-8)
	LevelChip     = slog.LevelDebug
	LevelBark     = slog.LevelInfo
	LevelKnot     = slog.LevelWarn
	LevelSplinter = slog.LevelError
	LevelTimber   = slog.Level(12)
)

func levelName(level slog.Level) string {
	switch level {
	case LevelSawdust:
		return "SAWDUST"
	case LevelChip:
		return "CHIP"
	case LevelBark:
		return "BARK"
	case LevelKnot:
		return "KNOT"
	case LevelSplinter:
		return "SPLINTER"
	case LevelTimber:
		return "TIMBER"
	default:
		return level.String()
	}
}
