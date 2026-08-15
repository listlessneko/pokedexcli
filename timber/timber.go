package timber

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

func New(writer io.Writer, level slog.Level) *Handler {
	return &Handler{
		Writer: writer,
		Level:  level,
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Level
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	time := record.Time.Format(time.RFC822)
	level := record.Level
	message := record.Message
	_, err := fmt.Fprintf(h.Writer, "(%s) [%s]: %s", time, level, message)
	if err != nil {
		return err
	}
	record.Attrs(
		func(a slog.Attr) bool {
			fmt.Fprintf(h.Writer, " %s=%v", a.Key, a.Value)
			return true
		})
	fmt.Fprintf(h.Writer, "\n")
	return nil
}

func (h *Handler) WithAttrs(attr []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(string string) slog.Handler {
	return h
}
