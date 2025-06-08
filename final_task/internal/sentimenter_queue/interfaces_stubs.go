package sentimeter_queue

import (
	"context"
	"math/rand/v2"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"pavlov061356/go-masters-homework/final_task/internal/storage"
	"time"
)

var _ storage.Interface = inMemory{}

var _ Sentimenter = sentimenterStub{}

type inMemory struct {
	outCH chan<- []models.Review
}

func newInMemory(outCh chan<- []models.Review) *inMemory {
	return &inMemory{outCH: outCh}
}

func (inMemory) NewUser(context.Context, models.User) (int, error) {
	return 0, nil
}
func (inMemory) GetUser(context.Context, int) (models.User, error) {
	return models.User{}, nil
}
func (inMemory) GetUserByEmail(context.Context, string) (models.User, error) {
	return models.User{}, nil
}
func (inMemory) UpdateUser(context.Context, models.User) error {
	return nil
}
func (inMemory) DeleteUser(context.Context, int) error {
	return nil
}

// CRUD методы для работы с услугами
func (inMemory) NewService(context.Context, models.Service) (int, error) {
	return 0, nil
}
func (inMemory) GetService(context.Context, int) (models.Service, error) {
	return models.Service{}, nil
}
func (inMemory) GetServices(context.Context) ([]models.Service, error) {
	return nil, nil
}
func (inMemory) UpdateService(context.Context, models.Service) error {
	return nil
}
func (inMemory) DeleteService(context.Context, int) error {
	return nil
}

// CRUD методы для работы с отзывами
func (inMemory) NewReview(context.Context, models.Review) (int, error) {
	return 0, nil
}
func (inMemory) GetReview(context.Context, int) (models.Review, error) {
	return models.Review{}, nil
}
func (inMemory) GetReviewsByService(context.Context, int) ([]models.Review, error) {
	return nil, nil
}
func (inMemory) GetReviewsByUser(context.Context, int) ([]models.Review, error) {
	return nil, nil
}
func (inMemory) UpdateReview(context.Context, models.Review) error {
	return nil
}
func (inMemory) DeleteReview(context.Context, models.Review) error {
	return nil
}
func (s inMemory) BatchUpdateReviewsSentiment(_ context.Context, reveiws []models.Review) error {
	s.outCH <- reveiws
	return nil
}
func (inMemory) RecomputeServicesScore(context.Context) error {
	return nil
}
func (inMemory) GetLastRecomputeTime(context.Context) (time.Time, error) {
	return time.Now(), nil
}
func (inMemory) GetUnsentimentedReviews(context.Context) ([]models.Review, error) {
	return nil, nil
}

type sentimenterStub struct {
}

func (sentimenterStub) GetReviewSentiment(ctx context.Context, review models.Review) (int, error) {
	return rand.IntN(3) + 1, nil
}
