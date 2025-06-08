package main

import (
	"context"
	"os"
	"os/signal"
	"pavlov061356/go-masters-homework/final_task/internal/server"
	"syscall"

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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	server := server.New(ctx)

	if err := server.Run(ctx); err != nil {
		log.Err(err).Msg("Ошибка запуска приложения")
		return
	}
}
