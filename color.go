package spretty

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	dim   = "\033[2m"

	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"

	cyan = "\033[36m"
	gray = "\033[90m"
)

func levelColor(level string) string {
	switch level {
	case "DEBUG":
		return blue
	case "INFO":
		return green
	case "WARN":
		return yellow
	case "ERROR":
		return red
	default:
		return gray
	}
}

func colorize(text, color string, noColor bool) string {
	if noColor {
		return text
	}
	return color + text + reset
}
