package spretty_test

import (
	"testing"

	spretty "github.com/mickamy/slog-pretty"
)

func TestLevelColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level string
		want  string
	}{
		{name: "DEBUG", level: "DEBUG", want: "\033[34m"},
		{name: "INFO", level: "INFO", want: "\033[32m"},
		{name: "WARN", level: "WARN", want: "\033[33m"},
		{name: "ERROR", level: "ERROR", want: "\033[31m"},
		{name: "unknown level", level: "TRACE", want: "\033[90m"},
		{name: "empty", level: "", want: "\033[90m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_ = t.Context()

			got := spretty.LevelColor(tt.level)
			if got != tt.want {
				t.Errorf("LevelColor(%q) = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		color   string
		noColor bool
		want    string
	}{
		{
			name:  "with color",
			text:  "INFO",
			color: "\033[32m",
			want:  "\033[32mINFO\033[0m",
		},
		{
			name:    "no color",
			text:    "INFO",
			color:   "\033[32m",
			noColor: true,
			want:    "INFO",
		},
		{
			name:  "empty text",
			text:  "",
			color: "\033[31m",
			want:  "\033[31m\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_ = t.Context()

			got := spretty.Colorize(tt.text, tt.color, tt.noColor)
			if got != tt.want {
				t.Errorf("Colorize(%q, %q, %v) = %q, want %q", tt.text, tt.color, tt.noColor, got, tt.want)
			}
		})
	}
}
