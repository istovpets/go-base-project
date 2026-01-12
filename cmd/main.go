package main

import (
	"base-project/internal/app"
	"log/slog"
)

func main() {
	a := app.New()
	defer a.Recover()

	slog.Debug("starting application...")

	a.Run()
	a.Stop()
}
