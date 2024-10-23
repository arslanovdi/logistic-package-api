// Package logger provides logging functionality
package logger

import (
	"log/slog"
	"os"
)

var options *slog.HandlerOptions
var loglevel *slog.LevelVar

// InitializeLogger initializes the slog logger
func InitializeLogger() {

	HidePassword := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "password" {
			return slog.String("password", "********")
		}
		return a
	}
	loglevel = &slog.LevelVar{}
	loglevel.Set(slog.LevelDebug)

	options = &slog.HandlerOptions{
		AddSource:   false,
		ReplaceAttr: HidePassword,
		Level:       loglevel,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, options))

	slog.SetDefault(logger)
	slog.Info("InitializeLogger", slog.String("level", loglevel.String()))
}

// SetLogLevel sets the level of the logger
// В пакете slog нет установки уровня логирования в дефолтном логере
func SetLogLevel(level slog.Level) {

	log := slog.With("func", "SetLogLevel")

	if options == nil {
		InitializeLogger()
	}

	loglevel.Set(level)
	log.Info("SetLogLevel", slog.String("level", level.String()))
}
