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
	parsedUrlsMux := sync.Mutex{}
	parsedUrls := map[string]struct{}{}

	var wg sync.WaitGroup

	crawlRequestMux := sync.Mutex{}
	crawlRequests := []crawlRequest{
		{
			url:   url,
			depth: 1,
		},
	}

	numWorkers := runtime.NumCPU()
	for i := range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					crawlRequestMux.Lock()
					if len(crawlRequests) == 0 {
						crawlRequestMux.Unlock()
						return
					}
					currentRequest := crawlRequests[0]
					crawlRequests = crawlRequests[1:]
					crawlRequestMux.Unlock()

					if p.maxDepth == 0 || (p.maxDepth != 0 && currentRequest.depth <= p.maxDepth) {
						log.Info().Msgf("Worker %d; Crawling url: %v", i, currentRequest)
						parserCtx, cancel := context.WithTimeout(ctx, p.timeout)
						parsed, err := p.parser.Parse(parserCtx, currentRequest.url)
						cancel()
						if err != nil {
							log.Error().Msgf("Failed to parse url: %v", err)
							continue
						}

						var tasks []string
						parsedUrlsMux.Lock()
						for _, url := range parsed {
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
					} else {
						continue
					}
				}
			}
		}()
	}

	parsingComplete := make(chan struct{})
	go func() {
		for parsed := range parsedUrlsChan {
			for _, url := range parsed {
				parsedUrlsMux.Lock()
				if _, ok := parsedUrls[url]; !ok {
					parsedUrls[url] = struct{}{}
				}
				parsedUrlsMux.Unlock()
			}
		}
		close(parsingComplete)
	}()

	wg.Wait()
	close(parsedUrlsChan)

	<-parsingComplete

	parsedUrlsMux.Lock()
	urls = make([]string, len(parsedUrls))
	for k := range parsedUrls {
		urls = append(urls, k)
	}
	parsedUrlsMux.Unlock()
	return
}
