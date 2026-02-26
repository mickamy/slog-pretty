package spretty

const (
	defaultTimeFormat = "15:04:05.000"
	defaultLevelWidth = 5
	defaultIndent     = "  "
)

type config struct {
	timeFormat string
	noColor    bool
	ignoreKeys map[string]struct{}
	levelWidth int
	indent     string
}

func defaultConfig(opts []Option) config {
	c := config{
		timeFormat: defaultTimeFormat,
		levelWidth: defaultLevelWidth,
		indent:     defaultIndent,
	}
	for _, o := range opts {
		o(&c)
	}
	return c
}

// Option configures formatting behavior.
type Option func(*config)

// WithTimeFormat sets the time format string (Go time package layout).
func WithTimeFormat(format string) Option {
	return func(c *config) {
		c.timeFormat = format
	}
}

// WithNoColor disables ANSI color output.
func WithNoColor() Option {
	return func(c *config) {
		c.noColor = true
	}
}

// WithIgnoreKeys specifies keys to omit from the formatted output.
func WithIgnoreKeys(keys ...string) Option {
	return func(c *config) {
		if c.ignoreKeys == nil {
			c.ignoreKeys = make(map[string]struct{}, len(keys))
		}
		for _, k := range keys {
			c.ignoreKeys[k] = struct{}{}
		}
	}
}
