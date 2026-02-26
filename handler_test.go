package spretty_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	spretty "github.com/mickamy/slog-pretty"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		hopts    *spretty.HandlerOptions
		opts     []spretty.Option
		log      func(l *slog.Logger)
		contains []string
		excludes []string
	}{
		{
			name:  "basic INFO",
			hopts: nil,
			opts:  []spretty.Option{spretty.WithNoColor()},
			log: func(l *slog.Logger) {
				l.Info("hello")
			},
			contains: []string{"INFO", "hello"},
		},
		{
			name:  "with attributes",
			hopts: nil,
			opts:  []spretty.Option{spretty.WithNoColor()},
			log: func(l *slog.Logger) {
				l.Info("started", "port", 8080, "host", "localhost")
			},
			contains: []string{"port=8080", "host=localhost"},
		},
		{
			name:  "level filtering",
			hopts: &spretty.HandlerOptions{Level: slog.LevelWarn},
			opts:  []spretty.Option{spretty.WithNoColor()},
			log: func(l *slog.Logger) {
				l.Info("hidden")
				l.Warn("visible")
			},
			contains: []string{"visible"},
			excludes: []string{"hidden"},
		},
		{
			name:  "with group",
			hopts: nil,
			opts:  []spretty.Option{spretty.WithNoColor()},
			log: func(l *slog.Logger) {
				l.WithGroup("req").Info("handled", "method", "GET")
			},
			contains: []string{"req.method=GET"},
		},
		{
			name:  "WithAttrs persists",
			hopts: nil,
			opts:  []spretty.Option{spretty.WithNoColor()},
			log: func(l *slog.Logger) {
				l.With("request_id", "abc-123").Info("done")
			},
			contains: []string{"request_id=abc-123", "done"},
		},
		{
			name:  "ERROR level",
			hopts: nil,
			opts:  []spretty.Option{spretty.WithNoColor()},
			log: func(l *slog.Logger) {
				l.Error("something failed", "err", "timeout")
			},
			contains: []string{"ERROR", "something failed", "err=timeout"},
		},
		{
			name:  "colored output",
			hopts: nil,
			opts:  []spretty.Option{},
			log: func(l *slog.Logger) {
				l.Info("hello")
			},
			contains: []string{"\033["},
		},
		{
			name:  "group attr",
			hopts: nil,
			opts:  []spretty.Option{spretty.WithNoColor()},
			log: func(l *slog.Logger) {
				l.Info("req",
					slog.Group("params",
						slog.String("table", "users"),
						slog.Int("limit", 100),
					),
				)
			},
			contains: []string{"params=", "table=users"},
		},
		{
			name:  "ignore keys",
			hopts: nil,
			opts: []spretty.Option{
				spretty.WithNoColor(),
				spretty.WithIgnoreKeys("secret"),
			},
			log: func(l *slog.Logger) {
				l.Info("test", "visible", "yes", "secret", "no")
			},
			contains: []string{"visible=yes"},
			excludes: []string{"secret"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_ = t.Context()

			var buf bytes.Buffer
			h := spretty.NewHandler(&buf, tt.hopts, tt.opts...)
			l := slog.New(h)

			tt.log(l)

			got := buf.String()

			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("output missing %q\ngot: %q", s, got)
				}
			}
			for _, s := range tt.excludes {
				if strings.Contains(got, s) {
					t.Errorf("output should not contain %q\ngot: %q", s, got)
				}
			}
		})
	}
}

func TestHandler_Enabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level slog.Leveler
		input slog.Level
		want  bool
	}{
		{name: "INFO enabled at INFO", level: slog.LevelInfo, input: slog.LevelInfo, want: true},
		{name: "WARN enabled at INFO", level: slog.LevelInfo, input: slog.LevelWarn, want: true},
		{name: "DEBUG disabled at INFO", level: slog.LevelInfo, input: slog.LevelDebug, want: false},
		{name: "INFO disabled at WARN", level: slog.LevelWarn, input: slog.LevelInfo, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			var buf bytes.Buffer
			h := spretty.NewHandler(&buf, &spretty.HandlerOptions{Level: tt.level})

			got := h.Enabled(ctx, tt.input)
			if got != tt.want {
				t.Errorf("Enabled(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
