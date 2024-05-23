package logging

import (
	"context"
	"log/slog"
	"os"
)

const (
	LevelTrace   = slog.Level(-8)
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelNotice  = slog.Level(2)
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12)
)

type LoggerConfig struct {
	LogLevel       string
	CommandHandler string
	AddSource      bool
}

type Logger struct {
	*slog.Logger
}

func NewLogger(logConfig LoggerConfig) (*Logger, error) {
	handlerOptions := &slog.HandlerOptions{
		Level:       slog.Level(-8),
		AddSource:   logConfig.AddSource,
		ReplaceAttr: replaceAttr,
	}
	var handler slog.Handler
	switch logConfig.CommandHandler {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	default:
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	}
	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
	}, nil
}

func (log *Logger) Fatal(msg string, args ...interface{}) {
	log.Logger.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

func (log *Logger) FatalCtx(ctx context.Context, msg string, args ...interface{}) {
	log.Logger.Log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)

		switch {
		case level < LevelDebug:
			a.Value = slog.StringValue("TRACE")
		case level < LevelInfo:
			a.Value = slog.StringValue("DEBUG")
		case level < LevelNotice:
			a.Value = slog.StringValue("INFO")
		case level < LevelWarning:
			a.Value = slog.StringValue("NOTICE")
		case level < LevelError:
			a.Value = slog.StringValue("WARNING")
		case level < LevelFatal:
			a.Value = slog.StringValue("ERROR")
		default:
			a.Value = slog.StringValue("FATAL")
		}
	}

	return a
}
