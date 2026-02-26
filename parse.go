package spretty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Parse attempts to parse a JSON line into a Record.
// Returns (record, true) on success, or (nil, false) if the line is not valid
// JSON or does not contain the expected slog fields.
func Parse(line []byte) (*Record, bool) {
	if len(line) == 0 || line[0] != '{' {
		return nil, false
	}

	dec := json.NewDecoder(bytes.NewReader(line))
	dec.UseNumber()

	tok, err := dec.Token()
	if err != nil {
		return nil, false
	}
	if _, ok := tok.(json.Delim); !ok {
		return nil, false
	}

	var rec Record
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return nil, false
		}
		key, ok := keyTok.(string)
		if !ok {
			return nil, false
		}

		switch key {
		case "time":
			rec.Time, err = parseTime(dec)
			if err != nil {
				return nil, false
			}
		case "level":
			rec.Level, ok = decodeString(dec)
			if !ok {
				return nil, false
			}
		case "msg":
			rec.Message, ok = decodeString(dec)
			if !ok {
				return nil, false
			}
		case "source":
			rec.Source, err = parseSource(dec)
			if err != nil {
				return nil, false
			}
		default:
			val, err := decodeValue(dec)
			if err != nil {
				return nil, false
			}
			rec.Attrs = append(rec.Attrs, Attr{Key: key, Value: val})
		}
	}

	if rec.Level == "" && rec.Message == "" {
		return nil, false
	}

	return &rec, true
}

func parseTime(dec *json.Decoder) (time.Time, error) {
	var s string
	if err := dec.Decode(&s); err != nil {
		return time.Time{}, fmt.Errorf("decoding time: %w", err)
	}

	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", s)
}

func parseSource(dec *json.Decoder) (*Source, error) {
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding source: %w", err)
	}

	var src Source
	if err := json.Unmarshal(raw, &src); err != nil {
		return nil, fmt.Errorf("unmarshalling source: %w", err)
	}
	return &src, nil
}

func decodeString(dec *json.Decoder) (string, bool) {
	var s string
	if err := dec.Decode(&s); err != nil {
		return "", false
	}
	return s, true
}

func decodeValue(dec *json.Decoder) (any, error) {
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding value: %w", err)
	}

	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()

	var v any
	if err := d.Decode(&v); err != nil {
		return nil, fmt.Errorf("unmarshalling value: %w", err)
	}
	return v, nil
}
