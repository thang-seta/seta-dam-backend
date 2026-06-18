package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	// Initialize slog to write JSON logs to stdout
	Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(Log)
}
