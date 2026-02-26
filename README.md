# slog-pretty

[![Sponsor](https://img.shields.io/badge/Sponsor-❤-ea4aaa?style=flat-square&logo=github)](https://github.com/sponsors/mickamy)

Pretty-printer for Go's `slog` JSON logs — inspired by [pino-pretty](https://github.com/pinojs/pino-pretty).

- Zero dependencies — standard library only
- Two ways to use: **CLI pipe** or **slog.Handler**
- Color-coded log levels, source locations, nested object expansion
- Respects `NO_COLOR` and auto-detects TTY

## Install

### CLI

```bash
# Homebrew
brew install mickamy/tap/spretty

# Go
go install github.com/mickamy/slog-pretty/cmd/spretty@latest
```

### Build from source

```bash
git clone https://github.com/mickamy/slog-pretty.git
cd slog-pretty
make install
```

### Docker

```dockerfile
ARG SPRETTY_VERSION=0.1.0
ADD https://github.com/mickamy/slog-pretty/releases/download/v${SPRETTY_VERSION}/slog-pretty_${SPRETTY_VERSION}_linux_${TARGETARCH}.tar.gz /tmp/spretty.tar.gz
RUN tar -xzf /tmp/spretty.tar.gz -C /usr/local/bin spretty && rm /tmp/spretty.tar.gz
```

Usage with `air` or `go run`:

```bash
air | spretty
go run ./cmd/server | spretty
```

### Library

```bash
go get github.com/mickamy/slog-pretty
```

## Quick Start

### As a CLI pipe

Pipe any JSON log output through `spretty`:

```bash
go run ./your-app | spretty
```

Output:

```
10:15:30.123 INFO  server started
  port=8080
  host=localhost
10:15:31.456 ERROR connection failed (main.connect /app/main.go:42)
  host=db.local
  retry=3
```

Non-JSON lines are passed through unchanged.

### As a slog.Handler

```go
package main

import (
	"log/slog"
	"os"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	logger := slog.New(spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	logger.Info("server started", "port", 8080)
	logger.With("request_id", "abc-123").Error("timeout", "host", "db.local")
}
```

## CLI Flags

| Flag              | Default        | Description                                             |
|-------------------|----------------|---------------------------------------------------------|
| `--time-format`   | `15:04:05.000` | Go [time format](https://pkg.go.dev/time#pkg-constants) |
| `--no-color`      | `false`        | Disable colored output                                  |
| `--ignore`        |                | Comma-separated keys to omit                            |
| `--version`, `-V` |                | Show version and exit                                   |

Colors are automatically disabled when stdout is not a TTY or when `NO_COLOR` is set.

## Handler Options

`NewHandler` accepts `HandlerOptions` and formatting `Option`s:

```go
spretty.NewHandler(w, hopts, opts...)
```

### HandlerOptions

| Field       | Type           | Default          | Description                  |
|-------------|----------------|------------------|------------------------------|
| `Level`     | `slog.Leveler` | `slog.LevelInfo` | Minimum log level            |
| `AddSource` | `bool`         | `false`          | Include source file and line |

### Formatting Options

| Function                  | Description                        |
|---------------------------|------------------------------------|
| `WithTimeFormat(format)`  | Set time format (Go layout string) |
| `WithNoColor()`           | Disable ANSI colors                |
| `WithIgnoreKeys(keys...)` | Omit specified keys from output    |

## Output Format

```
<time> <level> <message> (<source>)
  <key>=<value>
  <key>=
    <nested_key>=<value>
```

- **Time**: gray
- **Level**: DEBUG=blue, INFO=green, WARN=yellow, ERROR=red
- **Message**: bold
- **Source**: dim
- **Keys**: cyan
- Nested objects are indented; arrays are inline JSON

## License

[MIT](./LICENSE)
