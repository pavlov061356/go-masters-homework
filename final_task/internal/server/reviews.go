package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"strconv"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()

	review, err := parseBody[models.Review](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		span.RecordError(err)

		return
	}

	id, err := s.db.NewReview(r.Context(), review)
	if err != nil {
		log.Err(err).Msg("Ошибка при создании отзыва")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)

		return
	}

	log.Debug().Msgf("Создан новый отзыв с ID %d", id)
	span.AddEvent("Отзыв создан", trace.WithAttributes(attribute.Int("reviewID", id)))

	review.ID = id
	s.sentimeter.AddReview(review)
	log.Debug().Msgf("Отзыв с идентификатором %d отправлен в систему определния настроения отзыва", id)
	span.AddEvent("Отзыв отправлен в систему определния настроения отзыва", trace.WithAttributes(attribute.Int("reviewID", id)))

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(strconv.Itoa(id)))

	if err != nil {
		log.Err(err).Msg("Ошибка при записи ответа")
		span.RecordError(err)

		return
	}
}

func (s *Server) handleGetReview(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()

	reviewIDStr := r.PathValue("reviewID")
	if len(reviewIDStr) == 0 {
		http.Error(w, "не указан reviewID", http.StatusBadRequest)
		span.RecordError(fmt.Errorf("не указан reviewID"))

		return
	}

	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		http.Error(w, fmt.Errorf("ошибка при преобразовании reviewID в число: %w", err).Error(), http.StatusBadRequest)
		span.RecordError(fmt.Errorf("ошибка при преобразовании reviewID в число: %w", err))

		return
	}

	review, err := s.db.GetReview(r.Context(), reviewID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)

		return
	}

	err = json.NewEncoder(w).Encode(review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)

		return
	}

	span.SetStatus(codes.Ok, "OK")
}

func (s *Server) handleGetReviewsByService(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()

	serviceIDStr := r.PathValue("serviceID")
	if len(serviceIDStr) == 0 {
		http.Error(w, "не указан serviceID", http.StatusBadRequest)
		span.RecordError(fmt.Errorf("не указан serviceID"))

		return
	}

	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		http.Error(w, fmt.Errorf("ошибка при преобразовании serviceID в число: %w", err).Error(), http.StatusBadRequest)
		span.RecordError(fmt.Errorf("ошибка при преобразовании serviceID в число: %w", err))

		return
	}

	reviews, err := s.db.GetReviewsByService(r.Context(), serviceID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)

		return
	}

	err = json.NewEncoder(w).Encode(reviews)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)

		return
	}

	span.SetStatus(codes.Ok, "OK")
}
