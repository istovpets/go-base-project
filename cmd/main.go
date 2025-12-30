package main

import "base-project/internal/app"

func main() {
	a := app.New()
	defer a.Recover()

	log := a.Log()
	log.Debug("starting application...")

	a.Run()
	a.Stop()
}
