package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		},
	).With().Timestamp().Logger().With().Caller().Logger()
}

func main() {
	log.Debug().Msg("Начало работы приложения")
}
