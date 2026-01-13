package app

import "log/slog"

type ErrorPanic struct {
	Err error
}

func (e ErrorPanic) Error() string {
	return e.Err.Error()
}

func panicError(err error) {
	panic(ErrorPanic{Err: err})
}

func (a *App) Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(ErrorPanic); ok {
			slog.Error("failed to initialize application", slog.String("err", err.Error()))

			return
		}

		panic(r)
	}
}
