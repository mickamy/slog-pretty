package spretty_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	spretty "github.com/mickamy/slog-pretty"
)

func TestFormatter_Format(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 2, 26, 10, 15, 30, 123000000, time.UTC)

	tests := []struct {
		name     string
		opts     []spretty.Option
		record   spretty.Record
		contains []string
		excludes []string
	}{
		{
			name: "basic INFO",
			opts: []spretty.Option{spretty.WithNoColor()},
			record: spretty.Record{
				Time:    ts,
				Level:   "INFO",
				Message: "hello",
			},
			contains: []string{"10:15:30.123", "INFO ", "hello"},
		},
		{
			name: "ERROR with source",
			opts: []spretty.Option{spretty.WithNoColor()},
			record: spretty.Record{
				Time:    ts,
				Level:   "ERROR",
				Message: "fail",
				Source: &spretty.Source{
					Function: "main.run",
					File:     "/app/main.go",
					Line:     42,
				},
			},
			contains: []string{
				"ERROR", "fail",
				"(main.run /app/main.go:42)",
			},
		},
		{
			name: "with attributes",
			opts: []spretty.Option{spretty.WithNoColor()},
			record: spretty.Record{
				Time:    ts,
				Level:   "INFO",
				Message: "started",
				Attrs: []spretty.Attr{
					{Key: "port", Value: json.Number("8080")},
					{Key: "host", Value: "localhost"},
				},
			},
			contains: []string{
				"port=8080",
				"host=localhost",
			},
		},
		{
			name: "with nested object",
			opts: []spretty.Option{spretty.WithNoColor()},
			record: spretty.Record{
				Time:    ts,
				Level:   "INFO",
				Message: "req",
				Attrs: []spretty.Attr{
					{Key: "params", Value: map[string]any{
						"table": "users",
					}},
				},
			},
			contains: []string{
				"params=",
				"table=users",
			},
		},
		{
			name: "custom time format",
			opts: []spretty.Option{
				spretty.WithNoColor(),
				spretty.WithTimeFormat("2006-01-02"),
			},
			record: spretty.Record{
				Time:    ts,
				Level:   "INFO",
				Message: "test",
			},
			contains: []string{"2026-02-26"},
			excludes: []string{"10:15:30"},
		},
		{
			name: "ignore keys",
			opts: []spretty.Option{
				spretty.WithNoColor(),
				spretty.WithIgnoreKeys("secret"),
			},
			record: spretty.Record{
				Time:    ts,
				Level:   "INFO",
				Message: "test",
				Attrs: []spretty.Attr{
					{Key: "visible", Value: "yes"},
					{Key: "secret", Value: "hidden"},
				},
			},
			contains: []string{"visible=yes"},
			excludes: []string{"secret"},
		},
		{
			name: "colored output contains ANSI",
			opts: []spretty.Option{},
			record: spretty.Record{
				Time:    ts,
				Level:   "INFO",
				Message: "hello",
			},
			contains: []string{"\033["},
		},
		{
			name: "no extra attrs omits trailing newline",
			opts: []spretty.Option{spretty.WithNoColor()},
			record: spretty.Record{
				Time:    ts,
				Level:   "DEBUG",
				Message: "cache hit",
			},
			excludes: []string{"\n"},
		},
		{
			name: "boolean and null values",
			opts: []spretty.Option{spretty.WithNoColor()},
			record: spretty.Record{
				Level:   "INFO",
				Message: "test",
				Attrs: []spretty.Attr{
					{Key: "ok", Value: true},
					{Key: "err", Value: nil},
				},
			},
			contains: []string{"ok=true", "err=null"},
		},
		{
			name: "array value",
			opts: []spretty.Option{spretty.WithNoColor()},
			record: spretty.Record{
				Level:   "INFO",
				Message: "test",
				Attrs: []spretty.Attr{
					{Key: "ids", Value: []any{
						json.Number("1"),
						json.Number("2"),
					}},
				},
			},
			contains: []string{"ids=[1,2]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_ = t.Context()

			f := spretty.NewFormatter(tt.opts...)
			got := f.Format(&tt.record)

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
