package logger

import (
	"context"
	"log/slog"
	"os"
)

var logHandler *LevelHandler

type LevelHandler struct {
	level   slog.Level
	handler slog.Handler
}

// Enabled обработчик уровня логирования
func (lh *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= lh.level.Level()
}

func (lh *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return lh.handler.Handle(ctx, r)
}

func (lh *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewCustomLogger(lh.level, lh.handler.WithAttrs(attrs))
}

func (lh *LevelHandler) WithGroup(name string) slog.Handler {
	return NewCustomLogger(lh.level, lh.handler.WithGroup(name))
}

func (lh *LevelHandler) SetLogLevel(level slog.Level) {
	lh.level = level
}

func NewCustomLogger(level slog.Level, handler slog.Handler) *LevelHandler {
	return &LevelHandler{
		level:   level,
		handler: handler,
	}
}

func InitializeLogger() {

	HidePassword := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "password" {
			return slog.String("password", "********")
		}
		return a
	}

	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: HidePassword,
	})

	logHandler = NewCustomLogger(slog.LevelDebug, jsonH)

	logger := slog.New(logHandler)

	slog.SetDefault(logger)
}

// SetLogLevel sets the level of the logger
// В пакете slog нет установки уровня логирования в дефолтном логере
func SetLogLevel(level slog.Level) {
	logHandler.SetLogLevel(level)
}
