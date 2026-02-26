package spretty_test

import (
	"bytes"
	"strings"
	"testing"

	spretty "github.com/mickamy/slog-pretty"
)

func TestScanner_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		opts     []spretty.Option
		contains []string
		excludes []string
		lines    int
	}{
		{
			name: "single JSON line",
			input: `{"time":"2026-02-26T10:15:30Z","level":"INFO","msg":"hello","port":8080}
`,
			opts:     []spretty.Option{spretty.WithNoColor()},
			contains: []string{"10:15:30.000", "INFO", "hello", "port=8080"},
			lines:    1,
		},
		{
			name:     "non-JSON passthrough",
			input:    "plain text line\n",
			opts:     []spretty.Option{spretty.WithNoColor()},
			contains: []string{"plain text line"},
			lines:    1,
		},
		{
			name: "mixed JSON and non-JSON",
			input: `{"time":"2026-02-26T10:00:00Z","level":"INFO","msg":"first"}
not json
{"time":"2026-02-26T10:00:01Z","level":"ERROR","msg":"second"}
`,
			opts:     []spretty.Option{spretty.WithNoColor()},
			contains: []string{"INFO", "first", "not json", "ERROR", "second"},
			lines:    3,
		},
		{
			name:  "empty input",
			input: "",
			opts:  []spretty.Option{spretty.WithNoColor()},
			lines: 0,
		},
		{
			name: "with ignore keys",
			input: `{"time":"2026-02-26T10:00:00Z","level":"INFO","msg":"test","visible":"yes","secret":"no"}
`,
			opts: []spretty.Option{
				spretty.WithNoColor(),
				spretty.WithIgnoreKeys("secret"),
			},
			contains: []string{"visible=yes"},
			excludes: []string{"secret"},
			lines:    1,
		},
		{
			name: "multiple lines without trailing newline",
			input: `{"time":"2026-02-26T10:00:00Z","level":"INFO","msg":"a"}
{"time":"2026-02-26T10:00:01Z","level":"WARN","msg":"b"}`,
			opts:     []spretty.Option{spretty.WithNoColor()},
			contains: []string{"INFO", "a", "WARN", "b"},
			lines:    2,
		},
		{
			name:  "oversized line shows warning and continues",
			input: strings.Repeat("x", 2*1024*1024) + "\n" + `{"level":"INFO","msg":"after"}` + "\n",
			opts:  []spretty.Option{spretty.WithNoColor()},
			contains: []string{
				"WARN",
				"[spretty] line truncated",
				"max=1048576",
				"INFO",
				"after",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_ = t.Context()

			s := spretty.NewScanner(tt.opts...)
			var buf bytes.Buffer

			err := s.Scan(strings.NewReader(tt.input), &buf)
			if err != nil {
				t.Fatalf("Scan() error = %v", err)
			}

			got := buf.String()

			for _, c := range tt.contains {
				if !strings.Contains(got, c) {
					t.Errorf("output missing %q\ngot: %q", c, got)
				}
			}
			for _, e := range tt.excludes {
				if strings.Contains(got, e) {
					t.Errorf("output should not contain %q\ngot: %q", e, got)
				}
			}

			if tt.lines > 0 {
				gotLines := strings.Count(strings.TrimRight(got, "\n"), "\n") + 1
				if gotLines < tt.lines {
					t.Errorf("expected at least %d lines, got %d\noutput: %q",
						tt.lines, gotLines, got)
				}
			}
		})
	}
}
