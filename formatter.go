package spretty

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Formatter formats parsed Records into human-readable output.
type Formatter struct {
	cfg config
}

// NewFormatter creates a Formatter with the given options.
func NewFormatter(opts ...Option) *Formatter {
	return &Formatter{cfg: defaultConfig(opts)}
}

// Format returns the formatted representation of a Record.
func (f *Formatter) Format(r *Record) string {
	var b strings.Builder

	if !r.Time.IsZero() {
		b.WriteString(colorize(r.Time.Format(f.cfg.timeFormat), gray, f.cfg.noColor))
		b.WriteByte(' ')
	}

	if r.Level != "" {
		padded := fmt.Sprintf("%-*s", f.cfg.levelWidth, r.Level)
		b.WriteString(colorize(padded, levelColor(r.Level), f.cfg.noColor))
		b.WriteByte(' ')
	}

	b.WriteString(colorize(r.Message, bold, f.cfg.noColor))

	if r.Source != nil {
		src := fmt.Sprintf("(%s %s:%d)", r.Source.Function, r.Source.File, r.Source.Line)
		b.WriteByte(' ')
		b.WriteString(colorize(src, dim, f.cfg.noColor))
	}

	attrs := f.filterAttrs(r.Attrs)
	if len(attrs) > 0 {
		b.WriteByte('\n')
		f.writeAttrs(&b, attrs, f.cfg.indent)
	}

	return b.String()
}

func (f *Formatter) filterAttrs(attrs []Attr) []Attr {
	if len(f.cfg.ignoreKeys) == 0 {
		return attrs
	}
	filtered := make([]Attr, 0, len(attrs))
	for _, a := range attrs {
		if _, ignored := f.cfg.ignoreKeys[a.Key]; !ignored {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func (f *Formatter) writeAttrs(b *strings.Builder, attrs []Attr, prefix string) {
	for i, a := range attrs {
		b.WriteString(prefix)
		b.WriteString(colorize(a.Key, cyan, f.cfg.noColor))
		b.WriteString(colorize("=", gray, f.cfg.noColor))

		switch v := a.Value.(type) {
		case map[string]any:
			b.WriteByte('\n')
			f.writeMap(b, v, prefix+f.cfg.indent)
		default:
			b.WriteString(f.formatScalar(a.Value))
		}

		if i < len(attrs)-1 {
			b.WriteByte('\n')
		}
	}
}

func (f *Formatter) writeMap(b *strings.Builder, m map[string]any, prefix string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	for i, k := range keys {
		b.WriteString(prefix)
		b.WriteString(colorize(k, cyan, f.cfg.noColor))
		b.WriteString(colorize("=", gray, f.cfg.noColor))

		switch v := m[k].(type) {
		case map[string]any:
			b.WriteByte('\n')
			f.writeMap(b, v, prefix+f.cfg.indent)
		default:
			b.WriteString(f.formatScalar(m[k]))
		}

		if i < len(keys)-1 {
			b.WriteByte('\n')
		}
	}
}

func (f *Formatter) formatScalar(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	case []any:
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(data)
	default:
		return fmt.Sprintf("%v", v)
	}
}
