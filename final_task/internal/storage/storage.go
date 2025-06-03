package storage

import "pavlov061356/go-masters-homework/final_task/internal/models"

type Interface interface {
	// CRUD методы для работы с пользователями
	NewUser(models.User) (int, error)
	GetUser(int) (models.User, error)
	GetUserByEmail(string) (models.User, error)
	UpdateUser(models.User) error
	DeleteUser(int) error

	// CRUD методы для работы с услугами
	NewService(models.Service) (int, error)
	GetService(int) (models.Service, error)
	UpdateService(models.Service) error
	DeleteService(int) error

	// CRUD методы для работы с отзывами
	NewReview(models.Review) (int, error)
	GetReview(int) (models.Review, error)
	GetReviewsByService(int) ([]models.Review, error)
	GetReviewsByUser(int) ([]models.Review, error)
	UpdateReview(models.Review) error
	DeleteReview(int) error
}
