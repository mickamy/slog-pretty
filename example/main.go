package main

import (
	"log/slog"
	"os"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	handler := spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(handler)

	logger.Debug("cache hit", "key", "user:42")
	logger.Info("server started", "port", 8080, "host", "localhost")
	logger.Warn("slow query",
		slog.String("duration", "2.5s"),
		slog.Group("params", slog.String("table", "users"), slog.Int("limit", 100)),
	)
	logger.Error("connection failed", "err", "timeout", "host", "db.local", "retry", 3)

	logger.With("request_id", "abc-123").Info("request handled", "method", "GET", "path", "/api/users")
	logger.WithGroup("http").Info("response", "status", 200, "bytes", 1024)
}
