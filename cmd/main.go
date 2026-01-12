package main

import (
	"base-project/internal/app"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	a := app.New()
	defer a.Recover()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Debug("application starting...")
	err := a.Start(cancel)
	if err != nil {
		slog.Error("failed to start the application", slog.String("err", err.Error()))

		return
	}
	defer func() {
		err = a.Stop()
		if err != nil {
			slog.Error("Error occurred while stopping the application", slog.String("err", err.Error()))
		}

		slog.Debug("----------------------------------------------------------")
		slog.Debug("application stopped")
	}()

	err = a.CheckHealth(ctx)
	if err != nil {
		slog.Error("failed to check health", slog.String("err", err.Error()))
		cancel(err)

		return
	}

	slog.Debug("application started")
	slog.Debug("----------------------------------------------------------")

	err = a.Wait(ctx)
	if err != nil {
		slog.Error("shutdown due to error", slog.String("err", err.Error()))
	}

	slog.Debug("application stopping...")

}
