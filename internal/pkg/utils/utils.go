package utils

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

func LogErr(err error) slog.Attr {
	return slog.String("err", err.Error())
}

func ProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	dir := filepath.Dir(b)

	for range 10 {
		if dir == "/" || dir == "." {
			panic("project root not found")
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}

	panic("project root not found")
}
