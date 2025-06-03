package server

import (
	"net/http"
	"net/http/pprof"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"pavlov061356/go-masters-homework/final_task/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Sentimenter это интерфейс для определения настроения текста отзыва
type Sentimenter interface {
	GenerateSentiment(models.Review) error
}

type Server struct {
	router chi.Router

	db         storage.Interface
	sentimeter Sentimenter
}

func New(db storage.Interface, sentimeter Sentimenter) *Server {
	server := &Server{
		db:         db,
		sentimeter: sentimeter,
		router:     chi.NewRouter(),
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
	s.router.Get("/reviews", s.handleGetReviews)
	s.router.Get("/services/{serviceID}/score", s.handleGetServiceScore)

	s.router.Use(
		middleware.RequestID,
		middleware.RealIP,
	)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
