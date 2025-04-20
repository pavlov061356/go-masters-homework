package api

import (
	"net/http"
	"pavlov061356/go-masters-homework/2/internal/errs"
)

func (api *API) RandomErrHandler(w http.ResponseWriter, r *http.Request) {
	api.WriteError(w, r, errs.RandomErr())
}
