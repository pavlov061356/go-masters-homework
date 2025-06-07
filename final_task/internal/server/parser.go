package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
)

type Validator interface {
	Validate() error
}

// parseQuery - парсит модель запроса из URL query и валидирует его на корректность.
func parseQuery[T Validator](r *http.Request) (T, error) {
	var request T
	err := schema.NewDecoder().Decode(&request, r.URL.Query())
	if err != nil {
		return request, err
	}
	err = request.Validate()
	if err != nil {
		return request, err
	}

	return request, nil
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
