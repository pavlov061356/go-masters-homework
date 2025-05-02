package parser

import (
	"context"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    true,
			TimeFormat: "2006-01-02 15:04:05",
		},
	).With().Timestamp().Logger().With().Caller().Logger()
}

type Parser interface {
	Parse(ctx context.Context, url string) (urls []string, err error)
}

type URLParser struct {
	timeout  time.Duration
	maxDepth int
	parser   Parser
}

type crawlRequest struct {
	url   string
	depth int
}

func New(timeout time.Duration, maxDepth int, parser Parser) *URLParser {
	return &URLParser{
		timeout:  timeout,
		maxDepth: maxDepth,
		parser:   parser,
	}
}

func (p *URLParser) Parse(ctx context.Context, url string) (urls []string) {
	parsedUrlsChan := make(chan []string)

	var parsedUrlsMux sync.Mutex
	// Коллекция уникальных url, чтобы не парсить повторно ссылки
	parsedUrls := map[string]struct{}{}

	var wg sync.WaitGroup

	crawlRequestMux := sync.Mutex{}
	// Список ссылок для обработки
	crawlRequests := []crawlRequest{
		{
			url:   url,
			depth: 1,
		},
	}

	numWorkers := runtime.NumCPU()
	workerStatesMux := sync.Mutex{}
	// Состояние рабочих потоков, для определения когда все они завершили работу
	workerStates := make([]bool, numWorkers)
	for i := range numWorkers {
		workerStates[i] = false
		wg.Add(1)
		go func() {
			defer log.Info().Msgf("Worker %d; Exited", i)
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					crawlRequestMux.Lock()
					// Если нет ссылок для обработки, то проверяется статус всех остальных потоков,
					// если все работу не завершили, то ждём и переходим к следующему шагу цикла,
					// иначе завершаем работу этого потока
					if len(crawlRequests) == 0 {
						crawlRequestMux.Unlock()

						time.Sleep(time.Millisecond * 100)
						workerStatesMux.Lock()
						workerStates[i] = false

						var hasRunningWorkers bool
						for _, state := range workerStates {
							hasRunningWorkers = hasRunningWorkers || state
						}
						workerStatesMux.Unlock()

						if !hasRunningWorkers {
							return
						}
						continue
					}
					// Если есть ссылки для обработки, то берем первую из них и удаляем из списка
					currentRequest := crawlRequests[0]
					crawlRequests = crawlRequests[1:]
					crawlRequestMux.Unlock()

					// Если максимальная глубина обхода не достигнута, то начинаем обработку ссылки
					if p.maxDepth == 0 || (p.maxDepth != 0 && currentRequest.depth <= p.maxDepth) {

						// Помечаем рабочий поток как активный
						workerStatesMux.Lock()
						workerStates[i] = true
						workerStatesMux.Unlock()

						parserCtx, cancel := context.WithTimeout(ctx, p.timeout)
						parsed, err := p.parser.Parse(parserCtx, currentRequest.url)
						cancel()
						if err != nil {
							log.Error().Msgf("Failed to parse url: %v", err)
							continue
						}

						// tasks -- список ссылок для обработки, которые еще не были обработаны
						var tasks []string
						parsedUrlsMux.Lock()
						for _, url := range parsed {
							// Если ссылка не была обработана ранее, то добавляем ее в список ссылок для обработки
							if _, ok := parsedUrls[url]; !ok {
								crawlRequestMux.Lock()
								crawlRequests = append(crawlRequests, crawlRequest{
									url:   url,
									depth: currentRequest.depth + 1,
								})
								crawlRequestMux.Unlock()
								tasks = append(tasks, url)
							}
						}
						parsedUrlsMux.Unlock()
						if len(tasks) > 0 {
							parsedUrlsChan <- tasks
						}
					}
				}
			}
		}()
	}

	parsingComplete := make(chan struct{})
	go func() {
		for parsed := range parsedUrlsChan {
			parsedUrlsMux.Lock()
			for _, url := range parsed {
				parsedUrls[url] = struct{}{}
			}
			parsedUrlsMux.Unlock()
		}
		close(parsingComplete)
	}()

	wg.Wait()
	close(parsedUrlsChan)

	<-parsingComplete

	parsedUrlsMux.Lock()
	urls = make([]string, len(parsedUrls))
	i := 0
	for k := range parsedUrls {
		urls[i] = k
		i++
	}
	parsedUrlsMux.Unlock()
	return
}
