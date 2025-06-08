package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"strconv"

	"github.com/rs/zerolog/log"
)

func (s *Server) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	review, err := parseBody[models.Review](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.db.NewReview(r.Context(), review)
	if err != nil {
		log.Err(err).Msg("Ошибка при создании отзыва")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug().Msgf("Создан новый отзыв с ID %d", id)

	review.ID = id
	s.sentimeter.AddReview(review)
	log.Debug().Msgf("Отзыв с идентификатором %d отправлен в систему определния настроения отзыва", id)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(id)))
}

func (s *Server) handleGetReview(w http.ResponseWriter, r *http.Request) {
	reviewIDStr := r.PathValue("reviewID")
	if len(reviewIDStr) == 0 {
		http.Error(w, "не указан reviewID", http.StatusBadRequest)
		return
	}

	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		http.Error(w, fmt.Errorf("ошибка при преобразовании reviewID в число: %w", err).Error(), http.StatusBadRequest)
		return
	}

	review, err := s.db.GetReview(r.Context(), reviewID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(review)
}

func (s *Server) handleGetReviewsByService(w http.ResponseWriter, r *http.Request) {
	serviceIDStr := r.PathValue("serviceID")
	if len(serviceIDStr) == 0 {
		http.Error(w, "не указан serviceID", http.StatusBadRequest)
		return
	}

	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		http.Error(w, fmt.Errorf("ошибка при преобразовании serviceID в число: %w", err).Error(), http.StatusBadRequest)
		return
	}

	reviews, err := s.db.GetReviewsByService(r.Context(), serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}
