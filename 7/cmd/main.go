package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/valyala/fasthttp"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Запуск стандартного HTTP-сервера на порту 8080
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello from standard net/http server!")
		})

		log.Println("Standard HTTP server starting on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("Standard HTTP server failed: %v", err)
		}
	}()

	// Запуск fasthttp сервера на порту 8081
	go func() {
		defer wg.Done()
		handler := func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case "/":
				ctx.WriteString("Hello from fasthttp server!")
			default:
				ctx.Error("Not found", fasthttp.StatusNotFound)
			}
		}

		log.Println("FastHTTP server starting on :8081")
		if err := fasthttp.ListenAndServe(":8081", handler); err != nil {
			log.Fatalf("FastHTTP server failed: %v", err)
		}
	}()

	wg.Wait()
}
