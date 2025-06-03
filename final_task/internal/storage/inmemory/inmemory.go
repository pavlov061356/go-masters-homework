package inmemory

import (
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"pavlov061356/go-masters-homework/final_task/internal/storage"
)

var _ storage.Interface = &InMemory{}

type InMemory struct {
}

func (InMemory) NewUser(models.User) (int, error) {
	return 0, nil
}
func (InMemory) GetUser(int) (models.User, error) {
	return models.User{}, nil
}
func (InMemory) GetUserByEmail(string) (models.User, error) {
	return models.User{}, nil
}
func (InMemory) UpdateUser(models.User) error {
	return nil
}
func (InMemory) DeleteUser(int) error {
	return nil
}

// CRUD методы для работы с услугами
func (InMemory) NewService(models.Service) (int, error) {
	return 0, nil
}
func (InMemory) GetService(int) (models.Service, error) {
	return models.Service{}, nil
}
func (InMemory) GetServices() ([]models.Service, error) {
	return nil, nil
}
func (InMemory) UpdateService(models.Service) error {
	return nil
}
func (InMemory) DeleteService(int) error {
	return nil
}

// CRUD методы для работы с отзывами
func (InMemory) NewReview(models.Review) (int, error) {
	return 0, nil
}
func (InMemory) GetReview(int) (models.Review, error) {
	return models.Review{}, nil
}
func (InMemory) GetReviewsByService(int) ([]models.Review, error) {
	return nil, nil
}
func (InMemory) GetReviewsByUser(int) ([]models.Review, error) {
	return nil, nil
}
func (InMemory) UpdateReview(models.Review) error {
	return nil
}
func (InMemory) DeleteReview(int) error {
	return nil
}
