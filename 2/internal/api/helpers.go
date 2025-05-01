package api

import (
	"net/http"
	"pavlov061356/go-masters-homework/2/internal/errs"
)

func (*API) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var httpStatus int
	httpMessage := err.Error()

	switch err.(type) {
	case *errs.ErrBadRequest:
		httpStatus = http.StatusBadRequest
	case *errs.ErrNotFound:
		httpStatus = http.StatusNotFound
	case *errs.ErrConflict:
		httpStatus = http.StatusConflict
	default:
		httpStatus = http.StatusInternalServerError
	}

	w.WriteHeader(httpStatus)
	w.Write([]byte(httpMessage))
}
