package inmemory

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"pavlov061356/go-masters-homework/final_task/internal/storage"
)

var _ storage.Interface = &InMemory{}

type InMemory struct {
}

func (InMemory) NewUser(context.Context, models.User) (int, error) {
	return 0, nil
}
func (InMemory) GetUser(context.Context, int) (models.User, error) {
	return models.User{}, nil
}
func (InMemory) GetUserByEmail(context.Context, string) (models.User, error) {
	return models.User{}, nil
}
func (InMemory) UpdateUser(context.Context, models.User) error {
	return nil
}
func (InMemory) DeleteUser(context.Context, int) error {
	return nil
}

// CRUD методы для работы с услугами
func (InMemory) NewService(context.Context, models.Service) (int, error) {
	return 0, nil
}
func (InMemory) GetService(context.Context, int) (models.Service, error) {
	return models.Service{}, nil
}
func (InMemory) GetServices(context.Context) ([]models.Service, error) {
	return nil, nil
}
func (InMemory) UpdateService(context.Context, models.Service) error {
	return nil
}
func (InMemory) DeleteService(context.Context, int) error {
	return nil
}

// CRUD методы для работы с отзывами
func (InMemory) NewReview(context.Context, models.Review) (int, error) {
	return 0, nil
}
func (InMemory) GetReview(context.Context, int) (models.Review, error) {
	return models.Review{}, nil
}
func (InMemory) GetReviewsByService(context.Context, int) ([]models.Review, error) {
	return nil, nil
}
func (InMemory) GetReviewsByUser(context.Context, int) ([]models.Review, error) {
	return nil, nil
}
func (InMemory) UpdateReview(context.Context, models.Review) error {
	return nil
}
func (InMemory) DeleteReview(context.Context, models.Review) error {
	return nil
}
func (InMemory) BatchUpdateReviewsSentiment(context.Context, []models.Review) error {
	return nil
}
func (InMemory) RecomputeServicesScore(context.Context) error {
	return nil
}
