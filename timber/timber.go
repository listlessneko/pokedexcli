package timber

import (
	"context"
	"io"
	"log/slog"
)

func New(writer io.Writer, level slog.Level) *Handler {
	return &Handler{
		Writer: writer,
		Level:  level,
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (h *Handler) WithAttrs(attr []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(string string) slog.Handler {
	return h
}
