package timber

import (
	"io"
	"log/slog"
)

type Leveler = *slog.LevelVar

type Handler struct {
	Writer io.Writer
	Leveler  slog.Leveler
	Attrs  []slog.Attr
	Group  string
}

type Logger struct {
	logger *slog.Logger
}
