package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type AppLogger struct {
	Logger *zerolog.Logger
	LogLevel int
}


func NewLogger(logLevel int) *AppLogger {

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &AppLogger{
		Logger: &logger,
		LogLevel: logLevel,
	}
}

func (l *AppLogger) SetLevelInfo() {
	zerolog.SetGlobalLevel(zerolog.Level(l.LogLevel))
}