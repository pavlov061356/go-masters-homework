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
	parser   Parser
	timeout  time.Duration
	maxDepth int
}

func New(parser Parser, timeout time.Duration, maxDepth int) *URLParser {
	return &URLParser{
		parser:   parser,
		timeout:  timeout,
		maxDepth: maxDepth,
	}
}

func (c *URLParser) Parse(ctx context.Context, startURL string) []string {
	var (
		results []string
		visited sync.Map
		mu      sync.Mutex
		wg      sync.WaitGroup
		sema    = make(chan struct{}, runtime.NumCPU())
	)

	var parse func(context.Context, string, int)
	parse = func(ctx context.Context, currentURL string, depth int) {
		if c.maxDepth > 0 && depth > c.maxDepth {
			return
		}

		if _, loaded := visited.LoadOrStore(currentURL, struct{}{}); loaded {
			return
		}

		mu.Lock()
		results = append(results, currentURL)
		mu.Unlock()

		select {
		case <-ctx.Done():
			return
		case sema <- struct{}{}:
			defer func() { <-sema }()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
			defer cancel()

			urls, err := c.parser.Parse(reqCtx, currentURL)
			if err != nil {
				log.Error().Msgf("Failed to parse url: %v", err)
				return
			}

			for _, u := range urls {
				parse(ctx, u, depth+1)
			}
		}()
	}

	parse(ctx, startURL, 0)

	wg.Wait()

	return results
}
