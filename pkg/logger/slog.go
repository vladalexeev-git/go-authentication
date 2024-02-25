package logger

import (
	"log/slog"
	"os"
	"sso/config"
)

const (
	//format
	formatJson = "json"
	formatText = "text"
	//level
	levelDebug = "debug"
	levelInfo  = "info"
)

func SetupLogger(logCfg config.Logger) *slog.Logger {
	var opts = slog.HandlerOptions{
		Level:     slog.LevelInfo, //info level as a default
		AddSource: true,
	}

	switch logCfg.LogLevel {
	case levelDebug:
		opts.Level = slog.LevelDebug
	case levelInfo:
		opts.Level = slog.LevelInfo
	}

	if logCfg.Color {
		prettyOpts := PrettyHandlerOptions{
			SlogOpts: &opts,
		}
		handlerPretty := prettyOpts.NewPrettyHandler(os.Stdout, logCfg.Format)
		return slog.New(handlerPretty)
	}
	if logCfg.Format == formatText {
		handlerText := slog.NewTextHandler(os.Stdout, &opts)
		return slog.New(handlerText)
	}

	handlerJSON := slog.NewJSONHandler(os.Stdout, &opts)
	return slog.New(handlerJSON)
}
