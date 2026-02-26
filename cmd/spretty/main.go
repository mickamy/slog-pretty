package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	spretty "github.com/mickamy/slog-pretty"
)

var version = "dev"

func main() {
	fs := flag.NewFlagSet("spretty", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "spretty — pretty-print slog JSON output\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  <command> | spretty [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	timeFormat := fs.String("time-format", "15:04:05.000", "Go time format for timestamps")
	noColor := fs.Bool("no-color", false, "disable colored output")
	ignore := fs.String("ignore", "", "comma-separated keys to omit")
	showVersion := fs.Bool("version", false, "show version and exit")
	fs.BoolVar(showVersion, "V", false, "show version and exit (shorthand)")

	_ = fs.Parse(os.Args[1:])

	if *showVersion {
		fmt.Printf("spretty %s\n", version)
		return
	}

	var opts []spretty.Option
	opts = append(opts, spretty.WithTimeFormat(*timeFormat))

	if *noColor || os.Getenv("NO_COLOR") != "" || !isTerminal(os.Stdout) {
		opts = append(opts, spretty.WithNoColor())
	}

	if *ignore != "" {
		keys := strings.Split(*ignore, ",")
		opts = append(opts, spretty.WithIgnoreKeys(keys...))
	}

	s := spretty.NewScanner(opts...)
	if err := s.Scan(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "spretty: %v\n", err)
		os.Exit(1)
	}
}

func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
