package spretty

import "time"

// Record holds the parsed fields of a slog JSON log line.
type Record struct {
	Time    time.Time
	Level   string
	Message string
	Source  *Source
	Attrs   []Attr
}

// Source represents the slog source location.
type Source struct {
	Function string
	File     string
	Line     int
}

// Attr is an ordered key-value pair from a log line.
type Attr struct {
	Key   string
	Value any
}
