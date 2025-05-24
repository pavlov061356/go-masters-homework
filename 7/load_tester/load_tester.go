package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	urlStd     string
	urlFast    string
	requests   int
	concurrent int
)

func init() {
	flag.StringVar(&urlStd, "std", "http://localhost:8080", "URL стандартного HTTP сервера")
	flag.StringVar(&urlFast, "fast", "http://localhost:8081", "URL fasthttp сервера")
	flag.IntVar(&requests, "requests", 1000, "Общее количество запросов")
	flag.IntVar(&concurrent, "concurrent", 50, "Количество одновременных запросов")
	flag.Parse()
}

func main() {
	fmt.Printf("Нагрузочное тестирование:\n")
	fmt.Printf("Стандартный сервер: %s\n", urlStd)
	fmt.Printf("FastHTTP сервер:    %s\n", urlFast)
	fmt.Printf("Запросов: %d, Потоков: %d\n\n", requests, concurrent)

	// Тестируем стандартный сервер
	stdTime, stdSuccess := testServer(urlStd, "net/http", requests, concurrent)

	// Тестируем fasthttp сервер
	fastTime, fastSuccess := testServer(urlFast, "fasthttp", requests, concurrent)

	// Выводим результаты
	fmt.Printf("\nРезультаты:\n")
	fmt.Printf("Стандартный HTTP сервер:\n")
	fmt.Printf("  Время: %v\n", stdTime)
	fmt.Printf("  Успешных запросов: %d/%d (%.2f%%)\n",
		stdSuccess, requests, float64(stdSuccess)/float64(requests)*100)
	fmt.Printf("  Запросов в секунду: %.2f\n", float64(stdSuccess)/stdTime.Seconds())

	fmt.Printf("\nFastHTTP сервер:\n")
	fmt.Printf("  Время: %v\n", fastTime)
	fmt.Printf("  Успешных запросов: %d/%d (%.2f%%)\n",
		fastSuccess, requests, float64(fastSuccess)/float64(requests)*100)
	fmt.Printf("  Запросов в секунду: %.2f\n", float64(fastSuccess)/fastTime.Seconds())
}

func testServer(url, name string, totalRequests, concurrency int) (time.Duration, int64) {
	var success int64
	start := time.Now()

	// Создаем канал для ограничения количества одновременных запросов
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer func() {
				<-sem
				wg.Done()
			}()

			var resp *http.Response
			var err error

			if name == "fasthttp" {
				// Используем fasthttp клиент для fasthttp сервера
				statusCode, _, err := fasthttp.Get(nil, url)
				if err == nil && statusCode == http.StatusOK {
					atomic.AddInt64(&success, 1)
				}
			} else {
				// Используем стандартный http клиент
				resp, err = http.Get(url)
				if err == nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						atomic.AddInt64(&success, 1)
					}
				}
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("[%s] Завершено за %v, успешно: %d/%d (%.2f%%)\n",
		name, elapsed, success, totalRequests,
		float64(success)/float64(totalRequests)*100)

	return elapsed, success
}
