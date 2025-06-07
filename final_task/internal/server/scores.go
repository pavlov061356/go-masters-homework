package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) handleGetServiceScore(w http.ResponseWriter, r *http.Request) {
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

	service, err := s.db.GetService(r.Context(), serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(fmt.Appendf(nil, "%f", service.AvgScore))
}
