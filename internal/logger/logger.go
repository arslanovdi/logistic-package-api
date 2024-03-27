package logger

import (
	"log/slog"
	"os"
)

var options *slog.HandlerOptions
var loglevel *slog.LevelVar

func InitializeLogger() {

	HidePassword := func(groups []string, a slog.Attr) slog.Attr {
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, options))

	slog.SetDefault(logger)
	slog.Info("InitializeLogger", slog.String("level", loglevel.String()))
}

// SetLogLevel sets the level of the logger
// В пакете slog нет установки уровня логирования в дефолтном логере
func SetLogLevel(level slog.Level) {
	if options == nil {
		InitializeLogger()
	}

	loglevel.Set(level)
	slog.Info("SetLogLevel", slog.String("level", level.String()))
}
