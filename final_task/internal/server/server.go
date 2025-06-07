package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	sentimeter "pavlov061356/go-masters-homework/final_task/internal/sentimenter"
	"pavlov061356/go-masters-homework/final_task/internal/storage"
	"pavlov061356/go-masters-homework/final_task/internal/storage/postgres"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// Sentimenter это интерфейс для определения настроения текста отзыва
type Sentimenter interface {
	GenerateSentiment(models.Review) error
}

type Server struct {
	cfg *config.Config

	router chi.Router
	server *http.Server

	db         storage.Interface
	sentimeter Sentimenter
}

func New(ctx context.Context) *Server {
	router := chi.NewRouter()

	cfg, err := config.Load("")
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось загрузить конфигурацию")
	}

	db, err := postgres.New(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось создать подключение к базе данных")
	}

	sentimeter := sentimeter.New(ctx, cfg, db)

	server := &Server{
		db:         db,
		sentimeter: sentimeter,
		router:     router,
		cfg:        cfg,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%v", cfg.Port),
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  25 * time.Second,
		},
	}

	server.registerRoutes()

	return server
}

func (s *Server) registerRoutes() {
	// Эндпоинты pprof
	// http://localhost:8080/debug/pprof/
	s.router.Get("/debug/pprof/", pprof.Index)
	s.router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.Get("/debug/pprof/profile", pprof.Profile)
	s.router.Get("/debug/pprof/symbol", pprof.Symbol)
	s.router.Get("/debug/pprof/trace", pprof.Trace)
	s.router.Get("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	s.router.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	s.router.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.router.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	s.router.Get("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	s.router.Get("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	s.router.Get("/health", healthHandler)

	s.router.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Инициализация маршрутов
	s.router.Post("/reviews", s.handleCreateReview)
	s.router.Get("/reviews/{reviewID}", s.handleGetReview)
	s.router.Get("/services/{serviceID}/score", s.handleGetServiceScore)

	s.router.Use(
		middleware.RequestID,
		middleware.RealIP,
	)
}

// Run запускает сервер на порту cfg.Port
func (s *Server) Run(ctx context.Context) error {
	log.Info().Msgf("Запуск сервера на порту %d", s.cfg.Port)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info().Msg("Получен сигнал для остановки сервера")
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Ошибка при остановке сервера")
		}
	}()

	go s.refreshAvgScore(ctx)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// refreshAvgScore пересчитывает среднюю оценку услуг
func (s *Server) refreshAvgScore(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			lastRefreshTime, err := s.db.GetLastRecomputeTime(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Не удалось получить время последнего пересчёта средней оценки услуг")
				continue
			}

			if time.Since(lastRefreshTime) < time.Duration(s.cfg.AvgScoreRefreshTime) {
				time.Sleep(time.Duration(s.cfg.AvgScoreRefreshTime))
				continue
			}

			if err := s.db.RecomputeServicesScore(ctx); err != nil {
				log.Error().Err(err).Msg("Не удалось пересчитать среднюю оценку услуг")
			}
		}
	}
}
