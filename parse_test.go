package spretty_test

import (
	"testing"
	"time"

	spretty "github.com/mickamy/slog-pretty"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		wantOK bool
		check  func(t *testing.T, r *spretty.Record)
	}{
		{
			name:   "minimal slog line",
			input:  `{"time":"2026-02-26T10:15:30.123Z","level":"INFO","msg":"hello"}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()

				wantTime := time.Date(2026, 2, 26, 10, 15, 30, 123000000, time.UTC)
				if !r.Time.Equal(wantTime) {
					t.Errorf("Time = %v, want %v", r.Time, wantTime)
				}
				if r.Level != "INFO" {
					t.Errorf("Level = %q, want %q", r.Level, "INFO")
				}
				if r.Message != "hello" {
					t.Errorf("Message = %q, want %q", r.Message, "hello")
				}
				if r.Source != nil {
					t.Errorf("Source = %v, want nil", r.Source)
				}
				if len(r.Attrs) != 0 {
					t.Errorf("Attrs = %v, want empty", r.Attrs)
				}
			},
		},
		{
			name:   "with extra attributes",
			input:  `{"time":"2026-02-26T10:15:30Z","level":"INFO","msg":"hello","port":8080,"host":"localhost"}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()
				if len(r.Attrs) != 2 {
					t.Fatalf("Attrs length = %d, want 2", len(r.Attrs))
				}
				if r.Attrs[0].Key != "port" {
					t.Errorf("Attrs[0].Key = %q, want %q", r.Attrs[0].Key, "port")
				}
				if r.Attrs[1].Key != "host" {
					t.Errorf("Attrs[1].Key = %q, want %q", r.Attrs[1].Key, "host")
				}
			},
		},
		{
			name: "with source",
			input: `{"time":"2026-02-26T10:15:30Z","level":"ERROR","msg":"fail",` +
				`"source":{"function":"main.run","file":"/app/main.go","line":42}}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()
				if r.Source == nil {
					t.Fatal("Source is nil")
				}
				if r.Source.Function != "main.run" {
					t.Errorf("Source.Function = %q, want %q", r.Source.Function, "main.run")
				}
				if r.Source.File != "/app/main.go" {
					t.Errorf("Source.File = %q, want %q", r.Source.File, "/app/main.go")
				}
				if r.Source.Line != 42 {
					t.Errorf("Source.Line = %d, want %d", r.Source.Line, 42)
				}
			},
		},
		{
			name:   "with nested object",
			input:  `{"time":"2026-02-26T10:15:30Z","level":"INFO","msg":"req","params":{"table":"users","limit":100}}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()
				if len(r.Attrs) != 1 {
					t.Fatalf("Attrs length = %d, want 1", len(r.Attrs))
				}
				params, ok := r.Attrs[0].Value.(map[string]any)
				if !ok {
					t.Fatalf("Attrs[0].Value type = %T, want map[string]any", r.Attrs[0].Value)
				}
				if params["table"] != "users" {
					t.Errorf("params[table] = %v, want %q", params["table"], "users")
				}
			},
		},
		{
			name:   "RFC3339Nano timestamp",
			input:  `{"time":"2026-02-26T10:15:30.123456789Z","level":"DEBUG","msg":"tick"}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()
				if r.Time.Nanosecond() != 123456789 {
					t.Errorf("Nanosecond = %d, want 123456789", r.Time.Nanosecond())
				}
			},
		},
		{
			name:   "not JSON",
			input:  "plain text line",
			wantOK: false,
		},
		{
			name:   "empty line",
			input:  "",
			wantOK: false,
		},
		{
			name:   "empty object",
			input:  `{}`,
			wantOK: false,
		},
		{
			name:   "malformed JSON",
			input:  `{"time":`,
			wantOK: false,
		},
		{
			name:   "JSON array",
			input:  `[1,2,3]`,
			wantOK: false,
		},
		{
			name:   "level only is valid",
			input:  `{"level":"WARN"}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()
				if r.Level != "WARN" {
					t.Errorf("Level = %q, want %q", r.Level, "WARN")
				}
			},
		},
		{
			name:   "leading whitespace",
			input:  `  {"level":"INFO","msg":"hello"}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()
				if r.Message != "hello" {
					t.Errorf("Message = %q, want %q", r.Message, "hello")
				}
			},
		},
		{
			name:   "trailing garbage",
			input:  `{"level":"INFO","msg":"hello"} garbage`,
			wantOK: false,
		},
		{
			name:   "preserves attribute order",
			input:  `{"time":"2026-02-26T10:00:00Z","level":"INFO","msg":"test","z_key":"z","a_key":"a","m_key":"m"}`,
			wantOK: true,
			check: func(t *testing.T, r *spretty.Record) {
				t.Helper()
				wantKeys := []string{"z_key", "a_key", "m_key"}
				if len(r.Attrs) != len(wantKeys) {
					t.Fatalf("Attrs length = %d, want %d", len(r.Attrs), len(wantKeys))
				}
				for i, wk := range wantKeys {
					if r.Attrs[i].Key != wk {
						t.Errorf("Attrs[%d].Key = %q, want %q", i, r.Attrs[i].Key, wk)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec, ok := spretty.Parse([]byte(tt.input))
			if ok != tt.wantOK {
				t.Fatalf("Parse() ok = %v, want %v", ok, tt.wantOK)
			}
			if tt.check != nil && rec != nil {
				tt.check(t, rec)
			}
		})
	}
}
