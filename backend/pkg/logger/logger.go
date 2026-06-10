package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var log *slog.Logger

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	log = slog.New(handler)
}

func Info(msg string) {
	log.Info(msg)
}

func Infof(format string, args ...any) {
	log.Info(fmt.Sprintf(format, args...))
}

func Error(msg string) {
	log.Error(msg)
}

func Errorf(format string, args ...any) {
	log.Error(fmt.Sprintf(format, args...))
}
