package main

import (
	"context"
	"net/http"
	"pavlov061356/go-masters-homework/9/internal/telemetry"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	router := chi.NewRouter()
	router.Use(telemetry.TracingMiddleware)

	router.Get("/multiply/{number}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)
		defer span.End()

		number, err := strconv.Atoi(r.PathValue("number"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid number"))
			span.SetStatus(codes.Error, "invalid number")
			return
		}

		w.Write([]byte(strconv.Itoa(number * 2)))
	})

	telemetry.SetupOTelSDK(context.TODO(), "http://localhost:4318", "multiplier")

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Info().Msgf("multiplier up on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Err(err).Msg("Error running http server")
	}
}
