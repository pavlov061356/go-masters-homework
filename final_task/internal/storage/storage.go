package storage

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
)

type Interface interface {
	// CRUD методы для работы с пользователями
	NewUser(context.Context, models.User) (int, error)
	GetUser(context.Context, int) (models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
	UpdateUser(context.Context, models.User) error
	DeleteUser(context.Context, int) error

	// CRUD методы для работы с услугами
	NewService(context.Context, models.Service) (int, error)
	GetService(context.Context, int) (models.Service, error)
	UpdateService(context.Context, models.Service) error
	DeleteService(context.Context, int) error
	RecomputeServicesScore(context.Context) error

	// CRUD методы для работы с отзывами
	NewReview(context.Context, models.Review) (int, error)
	GetReview(context.Context, int) (models.Review, error)
	GetReviewsByService(context.Context, int) ([]models.Review, error)
	GetReviewsByUser(context.Context, int) ([]models.Review, error)
	UpdateReview(context.Context, models.Review) error
	BatchUpdateReviewsSentiment(context.Context, []models.Review) error
	DeleteReview(context.Context, models.Review) error
}
