package logger

import (
	"log/slog"
	"os"
)

func New(level slog.Level) (*slog.Logger, error) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler), nil
}
