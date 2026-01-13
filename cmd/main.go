package main

import (
	"base-project/internal/app"
	"base-project/internal/pkg/utils"
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
		slog.Error("failed to start the application", utils.LogErr(err))

		return
	}
	defer func() {
		err = a.Stop()
		if err != nil {
			slog.Error("error occurred while stopping the application", utils.LogErr(err))
		}

		slog.Debug("application stopped")
	}()

	err = a.CheckHealth(ctx)
	if err != nil {
		slog.Error("failed to check health", utils.LogErr(err))
		cancel(err)

		return
	}

	slog.Debug("application started")

	err = a.Wait(ctx)
	if err != nil {
		slog.Error("shutdown due to error", utils.LogErr(err))
	}

	slog.Debug("application stopping...")
}
