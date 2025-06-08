package sentimenter

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/models"
)

type Sentimenter struct {
	cfg config.Sentimenter
}

func New(cfg config.Sentimenter) *Sentimenter {
	return &Sentimenter{cfg: cfg}
}

func (s *Sentimenter) GetReviewSentiment(ctx context.Context, review models.Review) (int, error) {
	return 0, nil
}
