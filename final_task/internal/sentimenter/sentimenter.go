package sentimeter

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/config"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"pavlov061356/go-masters-homework/final_task/internal/storage"
)

type Sentimeter struct {
	ctx context.Context

	cfg *config.Config
	db  storage.Interface
}

func New(ctx context.Context, cfg *config.Config, db storage.Interface) *Sentimeter {
	return &Sentimeter{
		cfg: cfg,
		db:  db,
		ctx: ctx,
	}
}

func (s *Sentimeter) GenerateSentiment(models.Review) error {
	return nil
}
