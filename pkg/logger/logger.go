package logger

import (
	"os"
	"github.com/rs/zerolog"
)

type AppLogger struct {
	Logger *zerolog.Logger
	LogLevel int
}


func New(logLevel int) *AppLogger {

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &AppLogger{
		Logger: &logger,
		LogLevel: logLevel,
	}
}

func (l *AppLogger) SetLogLevel() {
	zerolog.SetGlobalLevel(zerolog.Level(l.LogLevel))
}