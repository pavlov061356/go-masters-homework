package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pavlov061356/go-masters-homework/2/internal/api"
	"strconv"
)

func main() {
	// Считываем переменную окружения PORT
	var port int
	portString := os.Getenv("PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		port = 8080
	}

	// Ргеистрируем хендлеры
	api := &api.API{}
	http.HandleFunc("/russian-roulette", api.RussianRouletteHandler)
	http.HandleFunc("/random-err", api.RandomErrHandler)

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
