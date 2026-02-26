package spretty

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"slices"
	"sync"
)

// Handler is a [slog.Handler] that writes human-readable, colorized log output.
type Handler struct {
	formatter *Formatter
	w         io.Writer
	mu        *sync.Mutex

	// preformatted attrs from WithAttrs / WithGroup calls.
	groups []string
	attrs  []Attr
}

// HandlerOptions holds configuration for [Handler].
type HandlerOptions struct {
	// Level reports the minimum level to log. Defaults to [slog.LevelInfo].
	Level slog.Leveler

	// AddSource causes the handler to include source location.
	AddSource bool
}

// NewHandler creates a [slog.Handler] that writes pretty-printed output to w.
// If hopts is nil, default options are used.
func NewHandler(w io.Writer, hopts *HandlerOptions, opts ...Option) *Handler {
	if hopts == nil {
		hopts = &HandlerOptions{}
	}
	if hopts.Level == nil {
		hopts.Level = slog.LevelInfo
	}

	cfg := newConfig(opts)
	cfg.handlerOpts = hopts

	return &Handler{
		formatter: &Formatter{cfg: cfg},
		w:         w,
		mu:        &sync.Mutex{},
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.formatter.cfg.handlerOpts.Level.Level()
}

// Handle formats and writes a log record.
func (h *Handler) Handle(_ context.Context, sr slog.Record) error {
	rec := &Record{
		Time:    sr.Time,
		Level:   sr.Level.String(),
		Message: sr.Message,
	}

	if h.formatter.cfg.handlerOpts.AddSource && sr.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{sr.PC})
		f, _ := fs.Next()
		rec.Source = &Source{
			Function: f.Function,
			File:     f.File,
			Line:     f.Line,
		}
	}

	// Prepend pre-formatted attrs from WithAttrs/WithGroup.
	rec.Attrs = append(rec.Attrs, h.attrs...)

	sr.Attrs(func(a slog.Attr) bool {
		rec.Attrs = append(rec.Attrs, slogAttrToAttr(h.groups, a))
		return true
	})

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := fmt.Fprintln(h.w, h.formatter.Format(rec))
	if err != nil {
		return fmt.Errorf("writing log: %w", err)
	}
	return nil
}

// WithAttrs returns a new Handler with the given attributes pre-formatted.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	converted := make([]Attr, 0, len(attrs))
	for _, a := range attrs {
		converted = append(converted, slogAttrToAttr(h.groups, a))
	}
	return &Handler{
		formatter: h.formatter,
		w:         h.w,
		mu:        h.mu,
		groups:    h.groups,
		attrs:     append(slices.Clone(h.attrs), converted...),
	}
}

// WithGroup returns a new Handler with the given group name.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &Handler{
		formatter: h.formatter,
		w:         h.w,
		mu:        h.mu,
		groups:    append(slices.Clone(h.groups), name),
		attrs:     slices.Clone(h.attrs),
	}
}

func slogAttrToAttr(groups []string, a slog.Attr) Attr {
	key := a.Key
	for i := len(groups) - 1; i >= 0; i-- {
		key = groups[i] + "." + key
	}

	if a.Value.Kind() == slog.KindGroup {
		m := make(map[string]any)
		for _, ga := range a.Value.Group() {
			m[ga.Key] = slogValueToAny(ga.Value)
		}
		return Attr{Key: key, Value: m}
	}

	return Attr{Key: key, Value: slogValueToAny(a.Value)}
}

func slogValueToAny(v slog.Value) any {
	switch v.Kind() {
	case slog.KindGroup:
		m := make(map[string]any)
		for _, a := range v.Group() {
			m[a.Key] = slogValueToAny(a.Value)
		}
		return m
	case slog.KindLogValuer:
		return slogValueToAny(v.Resolve())
	default:
		return v.Any()
	}
}
