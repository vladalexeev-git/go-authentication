package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"sso/config"
	"time"
)

func SetupLogger(logCfg config.Logger) *slog.Logger {
	var l *slog.Logger

	switch logCfg.Env {
	case "local":
		l = setupColorizedSlog()
	case "dev":
		l = setupColorizedSlog()
	case "prod":
		h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
		l = slog.New(h)
	default:
		h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
		l = slog.New(h)
	}
	return l
}

func setupColorizedSlog() *slog.Logger {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	)
	return logger
}
