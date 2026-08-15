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
	level := fmt.Sprintf("[%s]", record.Level)
	message := record.Message
	_, err := fmt.Fprintf(h.Writer, "(%s) %-7s: %s", time, level, message)
	if err != nil {
		return err
	}
	group := " "
	if h.Group != "" {
		group = " " + h.Group + "."
	}
	for _, a := range h.Attrs {
		fmt.Fprintf(h.Writer, "%s%s=%v", group, a.Key, a.Value)
	}
	record.Attrs(
		func(a slog.Attr) bool {
			fmt.Fprintf(h.Writer, "%s%s=%v", group, a.Key, a.Value)
			return true
		})
	fmt.Fprintf(h.Writer, "\n")
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.Attrs), len(h.Attrs)+len(attrs))
	copy(newAttrs, h.Attrs)
	newAttrs = append(newAttrs, attrs...)
	newH := New(h.Writer, h.Level)
	newH.Group = h.Group
	newH.Attrs = newAttrs
	return newH
}

func (h *Handler) WithGroup(name string) slog.Handler {
	newH := New(h.Writer, h.Level)
	newH.Attrs = h.Attrs
	if h.Group == "" {
		newH.Group = name
	} else {
		newH.Group = h.Group + "." + name
	}
	return newH
}
