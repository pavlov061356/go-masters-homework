package main

import (
	"fmt"
	"log"
	"net/http"
	"pavlov061356/go-masters-homework/2/internal/api"
)

func main() {
	port := 8080

	api := &api.API{}
	http.HandleFunc("/russian-roulette", api.RussianRouletteHandler)
	http.HandleFunc("/random-err", api.RandomErrHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
