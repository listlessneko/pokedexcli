package timber

import (
	"io"
	"log/slog"
)

type Handler struct {
	Writer io.Writer
	Level  slog.Level
	Attrs  []slog.Attr
	Group  string
}

type Logger struct {
	logger *slog.Logger
}
