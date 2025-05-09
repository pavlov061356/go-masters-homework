package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_fibonacciHandler(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantCode   int
		wantNumber int
	}{
		{
			name:     "/fibo?N=abc",
			url:      "?N=abc",
			wantCode: http.StatusBadRequest,
		},
		{
			name:       "/fibo?N=-1",
			url:        "?N=-1",
			wantCode:   http.StatusOK,
			wantNumber: 0,
		},
		{
			name:       "/fibo?N=0",
			url:        "?N=0",
			wantCode:   http.StatusOK,
			wantNumber: 0,
		},
		{
			name:       "/fibo?N=5",
			url:        "?N=5",
			wantCode:   http.StatusOK,
			wantNumber: 5,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(fibonacciHandler))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(server.URL + tt.url)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tt.wantCode {
				t.Errorf("got %d, want %d", res.StatusCode, tt.wantCode)
			}
			if res.StatusCode == http.StatusOK {
				var got int
				if _, err := fmt.Fscanf(res.Body, "%d", &got); err != nil {
					t.Fatal(err)
				}
				if got != tt.wantNumber {
					t.Errorf("got %d, want %d", got, tt.wantNumber)
				}
			}
		})
	}
}
