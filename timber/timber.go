package timber

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

func NewTimber() (*os.File, error) {
	err := os.MkdirAll("yard", 0755)
	if err != nil {
		return nil, err
	}
	timeLayout := "2006-01-02-150405"
	name := fmt.Sprintf("yard/timber-%s.log", time.Now().Format(timeLayout))
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func NewHandler(writer io.Writer, level slog.Level) *Handler {
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
	level := fmt.Sprintf("[%s]", levelName(record.Level))
	message := record.Message
	_, err := fmt.Fprintf(h.Writer, "(%s) %-10s: %s", time, level, message)
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
	newH := NewHandler(h.Writer, h.Level)
	newH.Group = h.Group
	newH.Attrs = newAttrs
	return newH
}

func (h *Handler) WithGroup(name string) slog.Handler {
	newH := NewHandler(h.Writer, h.Level)
	newH.Attrs = h.Attrs
	if h.Group == "" {
		newH.Group = name
	} else {
		newH.Group = h.Group + "." + name
	}
	return newH
}

func NewLogger(handler *Handler) *Logger {
	logger := slog.New(handler)
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Sawdust(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelSawdust, msg, args...)
}

func (l *Logger) Chip(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelChip, msg, args...)
}

func (l *Logger) Bark(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelBark, msg, args...)
}

func (l *Logger) Knot(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelKnot, msg, args...)
}

func (l *Logger) Splinter(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelSplinter, msg, args...)
}

func (l *Logger) Timber(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelTimber, msg, args...)
}

func (l *Logger) Branch(args ...any) *Logger {
	return &Logger{logger: l.logger.With(args...)}
}

func (l *Logger) BranchGroup(name string) *Logger {
	return &Logger{logger: l.logger.WithGroup(name)}
}
