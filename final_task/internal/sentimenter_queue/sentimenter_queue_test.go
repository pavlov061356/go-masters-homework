package sentimenterQueue

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"sync"
	"testing"
	"time"
)

func TestSentimeterQueue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.SentimenterQueue{
		MaxDBQueueLen:  100,
		MaxDBQueueWait: 10,
	}

	outCh := make(chan []models.Review)

	db := newInMemory(outCh)

	sentimenter := sentimenterStub{}

	queue, err := New(ctx, &cfg, db, sentimenter)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := range cfg.MaxDBQueueLen + 2 {
			queue.AddReview(models.Review{ID: i})
		}

		time.Sleep(time.Second * 20)
		cancel()
		close(outCh)
	}()

	var outReviews []models.Review
	go func() {
		defer wg.Done()
		for reviews := range outCh {
			outReviews = append(outReviews, reviews...)
		}
	}()

	wg.Wait()

	if len(outReviews) != cfg.MaxDBQueueLen+1 {
		t.Fatalf("неверное количество отзывов: %d", len(outReviews))
	}
}
