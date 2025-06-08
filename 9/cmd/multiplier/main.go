package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"pavlov061356/go-masters-homework/9/internal/telemetry"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

		result, err := multiplyNumber(r.Context(), number+2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			span.SetStatus(codes.Error, err.Error())
			return
		}

		w.Write([]byte(strconv.Itoa(result)))
	})

	telemetry.SetupOTelSDK(context.TODO(), "http://localhost:4318", "summer")

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Info().Msgf("summer up on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Err(err).Msg("Error running http server")
	}
}

func multiplyNumber(ctx context.Context, number int) (int, error) {
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/multiply/"+strconv.Itoa(number), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to build request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to request bar service: %v", err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read bar service response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bar service responded with status %d and body '%s'", resp.StatusCode, response)
	}

	result, err := strconv.Atoi(string(response))
	if err != nil {
		return 0, fmt.Errorf("failed to parse bar service response: %v", err)
	}

	return result, nil

}
