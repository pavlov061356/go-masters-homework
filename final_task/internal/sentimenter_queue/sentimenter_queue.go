package sentimeter_queue

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"pavlov061356/go-masters-homework/final_task/internal/storage"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Sentimenter interface {
	GetReviewSentiment(ctx context.Context, review models.Review) (int, error)
}

type SentimenterQueue struct {
	ctx context.Context

	cfg         *config.SentimenterQueue
	db          storage.Interface
	sentimenter Sentimenter

	mux          sync.Mutex
	reviewsQueue []models.Review
}

func New(ctx context.Context, cfg *config.SentimenterQueue, db storage.Interface, sentimenter Sentimenter) *SentimenterQueue {
	queue := SentimenterQueue{
		cfg:         cfg,
		db:          db,
		ctx:         ctx,
		sentimenter: sentimenter,
	}

	reviews, err := queue.db.GetUnsentimentedReviews(ctx)
	if err != nil {
		log.Err(err).Msg("Ошибка при получении отзывов без настроения отзыва")
		return nil
	}
	queue.reviewsQueue = reviews

	return &queue
}

func (s *SentimenterQueue) Run() {
	var wg sync.WaitGroup

	wg.Add(2)

	outCh := make(chan models.Review)

	go func() {
		defer wg.Done()
		for {
			s.mux.Lock()
			var review models.Review
			if len(s.reviewsQueue) > 0 {
				review = s.reviewsQueue[0]
				s.reviewsQueue = s.reviewsQueue[1:]
			}
			s.mux.Unlock()

			if review.ID == 0 {
				time.Sleep(time.Second)
				continue
			}

			sentiment, err := s.sentimenter.GetReviewSentiment(s.ctx, review)
			if err != nil {
				log.Err(err).Msg("Ошибка при получении настроения отзыва")
				s.mux.Lock()
				s.reviewsQueue = append(s.reviewsQueue, review)
				s.mux.Unlock()
				continue
			}

			review.Sentiment = sentiment

			select {
			case <-s.ctx.Done():
				close(outCh)
				return
			case outCh <- review:
			}
		}
	}()

	go func() {
		defer wg.Done()
		var outReviews []models.Review
		for {
			select {
			case <-s.ctx.Done():
				return
			case outReview := <-outCh:
				log.Debug().Int("len", len(outReviews)).Msg("Добавление отзыва в очередь на сохранение")
				outReviews = append(outReviews, outReview)
				if len(outReviews) >= s.cfg.MaxDBQueueLen {
					log.Debug().Int("len", len(outReviews)).Msg("Сохранение настроений отзывов по превышению максимальной длины очереди")
					err := s.db.BatchUpdateReviewsSentiment(s.ctx, outReviews)
					if err != nil {
						log.Err(err).Msg("Ошибка при сохранении настроений отзывов")
						continue
					}
					log.Debug().Msg("Сохранение настроений отзывов завершено")
					outReviews = nil
				}
			case <-time.After(time.Second * 10):
				log.Debug().Int("len", len(outReviews)).Msg("Сохранение настроений отзывов по таймауту")
				err := s.db.BatchUpdateReviewsSentiment(s.ctx, outReviews)
				if err != nil {
					log.Err(err).Msg("Ошибка при сохранении настроений отзывов")
					continue
				}
				log.Debug().Msg("Сохранение настроений отзывов завершено")
				outReviews = nil
			}
		}
	}()

	wg.Wait()
}

// AddReview добавляет отзыв в очередь на обработку.
func (s *SentimenterQueue) AddReview(review models.Review) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.reviewsQueue = append(s.reviewsQueue, review)
}
