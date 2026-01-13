package utils

import "log/slog"

func LogErr(err error) slog.Attr {
	return slog.String("err", err.Error())
}
