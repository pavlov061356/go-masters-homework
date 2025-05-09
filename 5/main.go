package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

func pointer[T any](v T) *T {
	return &v
}

func main() {
	http.HandleFunc("/ollama", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing 'q' query parameter"))
			return
		}

		addr, err := url.Parse("http://localhost:11434")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fmt.Appendf(nil, "Не удалось создать Ollama клиент: %v", err))
			return
		}
		client := api.NewClient(addr, http.DefaultClient)

		// model := "gemma3:1b"
		model := "llama3.2:1b"

		generateReq := &api.GenerateRequest{
			Model:  model,
			Prompt: query,
			Stream: pointer(true),
		}

		var response string
		generateRespFunc := func(resp api.GenerateResponse) error {
			if !resp.Done {
				response += resp.Response
			}
			return nil
		}

		err = client.Generate(r.Context(), generateReq, generateRespFunc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
