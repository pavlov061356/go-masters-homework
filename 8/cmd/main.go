package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/emersion/go-imap/server"
	"github.com/rs/zerolog/log"
)

func main() {
	// Создаем контекст для graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	// Инициализируем сервер
	srv := server.New()

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.Start(ctx); err != nil {
			log.Error().Err(err).Msg("Ошибка при запуске сервера")
			cancel()
		}
	}()

	<-ctx.Done()

	log.Info().Msg("Получен сигнал на завершение работы")
	cancel()

	// Ожидаем завершения работы сервера
	<-ctx.Done()
	log.Info().Msg("Сервер успешно остановлен")
}
