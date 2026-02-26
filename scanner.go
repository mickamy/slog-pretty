package spretty

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const maxLineSize = 1024 * 1024 // 1MB

// Scanner reads lines from an io.Reader, formats slog JSON lines,
// and writes the output to an io.Writer. Non-JSON lines are passed through.
// Lines exceeding 1MB are passed through with a truncation warning.
type Scanner struct {
	formatter *Formatter
}

// NewScanner creates a Scanner with the given options.
func NewScanner(opts ...Option) *Scanner {
	return &Scanner{formatter: NewFormatter(opts...)}
}

// Scan reads from r and writes formatted output to w.
// It processes input line-by-line and returns any I/O error encountered.
func (s *Scanner) Scan(r io.Reader, w io.Writer) error {
	br := bufio.NewReaderSize(r, maxLineSize)

	for {
		line, err := readLine(br)
		if len(line) > 0 {
			if writeErr := s.processLine(w, line); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("reading input: %w", err)
		}
	}
}

func readLine(br *bufio.Reader) ([]byte, error) {
	var buf []byte
	for {
		fragment, isPrefix, err := br.ReadLine()
		buf = append(buf, fragment...)

		if !isPrefix || err != nil {
			if err != nil {
				err = fmt.Errorf("reading line: %w", err)
			}
			return buf, err
		}

		if len(buf) > maxLineSize {
			// Discard the rest of the oversized line.
			for isPrefix && err == nil {
				_, isPrefix, err = br.ReadLine()
			}
			if err != nil {
				return buf, fmt.Errorf("discarding oversized line: %w", err)
			}
			return buf, nil
		}
	}
}

func (s *Scanner) processLine(w io.Writer, line []byte) error {
	if len(line) > maxLineSize {
		return s.writeOverflow(w, line)
	}

	rec, ok := Parse(line)
	if ok {
		_, err := fmt.Fprintln(w, s.formatter.Format(rec))
		if err != nil {
			return fmt.Errorf("writing formatted line: %w", err)
		}
		return nil
	}

	_, err := fmt.Fprintln(w, string(line))
	if err != nil {
		return fmt.Errorf("writing passthrough line: %w", err)
	}
	return nil
}

//nolint:gosec // output is log text, not user-facing HTML
func (s *Scanner) writeOverflow(w io.Writer, line []byte) error {
	rec := &Record{
		Level:   "WARN",
		Message: "[spretty] line truncated",
		Attrs: []Attr{
			{Key: "max", Value: json.Number(strconv.Itoa(maxLineSize))},
			{Key: "size", Value: json.Number(strconv.Itoa(len(line)))},
		},
	}
	_, err := fmt.Fprintln(w, s.formatter.Format(rec))
	if err != nil {
		return fmt.Errorf("writing overflow line: %w", err)
	}
	return nil
}
