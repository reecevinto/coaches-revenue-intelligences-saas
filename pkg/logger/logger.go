package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(env string) zerolog.Logger {
	if env == "development" {
		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger().
			Output(zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			})
	}

	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}