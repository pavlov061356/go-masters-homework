package api

import (
	"math/rand/v2"
	"net/http"
	"pavlov061356/go-masters-homework/2/internal/errs"
)

func (api *API) RussianRouletteHandler(w http.ResponseWriter, r *http.Request) {
	rnd := rand.IntN(2)
	if rnd == 0 {
		w.Write([]byte("Вы выжили!"))
		return
	} else {
		api.WriteError(w, r, errs.RandomErr())
		return
	}
}
