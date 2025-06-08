package server

import (
	"fmt"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) handleGetServiceScore(w http.ResponseWriter, r *http.Request) {
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

	service, err := s.db.GetService(r.Context(), serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)
		return
	}

	span.SetStatus(codes.Ok, "Запрос успешно обработан")
	w.Write(fmt.Appendf(nil, "%f", service.AvgScore))
}
