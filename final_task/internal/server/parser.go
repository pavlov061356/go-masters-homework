package server

import (
	"encoding/json"
	"net/http"
)

type Validator interface {
	Validate() error
}

// parseBody - парсит модель запроса из тела запроса и валидирует его на корректность.
func parseBody[T Validator](r *http.Request) (T, error) {
	var request T
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		return request, err
	}

	err = request.Validate()

	if err != nil {
		return request, err
	}

	return request, nil
}
