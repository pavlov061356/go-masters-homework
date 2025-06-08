package sentimenter

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/consts"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	sentimenterQueue "pavlov061356/go-masters-homework/final_task/internal/sentimenter_queue"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rs/zerolog/log"
)

var _ sentimenterQueue.Sentimenter = (*Sentimenter)(nil)

type Sentimenter struct {
	cfg config.Sentimenter

	client *api.Client
}

// New создает новый объект для работы с Ollama.
func New(cfg config.Sentimenter) (*Sentimenter, error) {
	addr, err := url.Parse(cfg.Addr)
	if err != nil {
		log.Err(err).Msgf("Не удалось создать Ollama клиент: %v", err)
		return nil, err
	}

	// http.DefaultClient.
	client := api.NewClient(addr, http.DefaultClient)
	return &Sentimenter{cfg: cfg, client: client}, nil
}

func sentimentFromReviewScore(score int) int {
	switch {
	case score <= 2:
		return consts.SentimentNegative
	case score <= 4:
		return consts.SentimentNeutral
	default:
		return consts.SentimentPositive
	}
}

func (s *Sentimenter) GetReviewSentiment(ctx context.Context, review models.Review) (int, error) {
	if review.Content == "" {
		return sentimentFromReviewScore(review.Score), nil
	}
	req := &api.GenerateRequest{
		Model:  s.cfg.Model,
		Prompt: fmt.Sprintf(`Generate sentiment of the review: "%s"; return sentiment as integer: 1 - positive; 2 - neutral; 3 - negative;`, review.Content),
		Stream: new(bool),
	}

	type response struct {
		sentiment int
		err       error
	}

	respChan := make(chan response)

	respFunc := func(resp api.GenerateResponse) error {
		log.Debug().Msg(resp.Response)
		sentiment, err := strconv.Atoi(strings.TrimSpace(resp.Response[:1]))
		var response response
		if err != nil {
			response.err = err
		}
		response.sentiment = sentiment
		respChan <- response
		close(respChan)
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Timeout)*time.Second)
	defer cancel()

	var gotResponse response

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for response := range respChan {
			gotResponse = response
		}
	}()

	err := s.client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Err(err).Msg("Ошибка получения результата запроса к Ollama")
	}

	wg.Wait()

	if gotResponse.err != nil {
		log.Err(gotResponse.err).Msg("Ошибка при обработке результата запроса к Ollama")
	}

	if gotResponse.sentiment == 0 {
		return sentimentFromReviewScore(review.Score), nil
	}
	return gotResponse.sentiment, nil
}
