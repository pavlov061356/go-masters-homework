package main

import (
	"log"
	"net/http"
	"strconv"
)

func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	num := r.URL.Query().Get("N")
	n, err := strconv.Atoi(num)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	result := fibonacci(n)
	w.Write([]byte(strconv.Itoa(result)))
}

func main() {
	http.HandleFunc("/fibo", fibonacciHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
