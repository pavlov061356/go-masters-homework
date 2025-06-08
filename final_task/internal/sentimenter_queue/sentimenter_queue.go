package sentimenterQueue

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/metrics"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"pavlov061356/go-masters-homework/final_task/internal/storage"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Sentimenter интерфейс сервиса для получения настроения отзыва.
type Sentimenter interface {
	GetReviewSentiment(ctx context.Context, review models.Review) (int, error)
}

// SentimenterQueue структура организации очереди получения настроения отзыва и записи полученных данных в БД.
type SentimenterQueue struct {
	cfg         *config.SentimenterQueue
	db          storage.Interface
	sentimenter Sentimenter

	mux          sync.Mutex
	reviewsQueue []models.Review
}

// New создает очередь получения настроения отзыва и записи полученных данных в БД.
func New(ctx context.Context, cfg *config.SentimenterQueue, db storage.Interface, sentimenter Sentimenter) *SentimenterQueue {
	queue := SentimenterQueue{
		cfg:         cfg,
		db:          db,
		sentimenter: sentimenter,
	}

	reviews, err := queue.db.GetUnsentimentedReviews(ctx)
	if err != nil {
		log.Err(err).Msg("Ошибка при получении отзывов без настроения отзыва")
		return nil
	}

	queue.reviewsQueue = reviews

	go queue.run(ctx)

	return &queue
}

// Run запускает очередь получения настроения отзыва и записи полученных данных в БД.
func (s *SentimenterQueue) run(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(3)

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

			sentiment, err := s.sentimenter.GetReviewSentiment(ctx, review)
			if err != nil {
				log.Err(err).Msg("Ошибка при получении настроения отзыва")
				s.mux.Lock()
				s.reviewsQueue = append(s.reviewsQueue, review)
				s.mux.Unlock()

				continue
			}

			metrics.ReviewsSentimentDistribution.WithLabelValues().Observe(float64(sentiment))
			review.Sentiment = sentiment

			select {
			case <-ctx.Done():
				close(outCh)
				return
			case outCh <- review:
			}
		}
	}()

	mux := &sync.Mutex{}

	var outReviews []models.Review

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute * 5):
			mux.Lock()
			metrics.SentimenterQueueLength.Add(float64(len(outReviews)))
			mux.Unlock()
		}
	}()

	// Принцип работы такой, что либо набирается максимальная длина очереди, либо по таймауту очередь отправлется на сохранение в БД.
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case outReview := <-outCh:
				func() {
					mux.Lock()
					defer mux.Unlock()

					log.Debug().Int("len", len(outReviews)).Msg("Добавление отзыва в очередь на сохранение")
					outReviews = append(outReviews, outReview)

					if len(outReviews) >= s.cfg.MaxDBQueueLen {
						log.Debug().Int("len", len(outReviews)).Msg("Сохранение настроений отзывов по превышению максимальной длины очереди")
						err := s.db.BatchUpdateReviewsSentiment(ctx, outReviews)

						if err != nil {
							log.Err(err).Msg("Ошибка при сохранении настроений отзывов")
							return
						}

						log.Debug().Msg("Сохранение настроений отзывов завершено")

						outReviews = nil
					}
				}()
			case <-time.After(time.Second * 10):
				func() {
					mux.Lock()
					defer mux.Unlock()

					log.Debug().Int("len", len(outReviews)).Msg("Сохранение настроений отзывов по таймауту")

					if len(outReviews) == 0 {
						return
					}

					err := s.db.BatchUpdateReviewsSentiment(ctx, outReviews)

					if err != nil {
						log.Err(err).Msg("Ошибка при сохранении настроений отзывов")
						return
					}

					log.Debug().Msg("Сохранение настроений отзывов завершено")

					outReviews = nil
				}()
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
