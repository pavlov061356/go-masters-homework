package storage

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"time"
)

type Interface interface {
	// CRUD методы для работы с пользователями

	// NewUser создаёт нового пользователя и возвращает его идентификатор.
	NewUser(context.Context, models.User) (int, error)
	// GetUser возвращает пользователя по идентификатору.
	GetUser(context.Context, int) (models.User, error)
	// GetUserByEmail возвращает пользователя по email.
	GetUserByEmail(context.Context, string) (models.User, error)
	// UpdateUser обновляет данные пользователя.
	UpdateUser(context.Context, models.User) error
	// DeleteUser удаляет пользователя по идентификатору.
	DeleteUser(context.Context, int) error

	// CRUD методы для работы с услугами

	// NewService создает новую услугу и возвращает ее идентификатор.
	NewService(context.Context, models.Service) (int, error)
	// GetService возвращает услугу по ее идентификатору.
	GetService(context.Context, int) (models.Service, error)
	// GetServices возвращает все услуги.
	GetServices(context.Context) ([]models.Service, error)
	// UpdateService обновляет информацию об услуге по ее идентификатору.
	UpdateService(context.Context, models.Service) error
	// DeleteService удаляет услугу по ее идентификатору.
	DeleteService(context.Context, int) error
	// RecomputeServicesScore пересчитывает средний рейтинг услуг.
	RecomputeServicesScore(context.Context) error
	// GetLastRecomputeTime возвращает время последнего пересчета рейтингов услуг.
	GetLastRecomputeTime(context.Context) (time.Time, error)

	// CRUD методы для работы с отзывами

	// NewReview метод для создания отзыва.
	// Обновляет среднюю оценку услуги и возвращает id созданного отзыва.
	NewReview(context.Context, models.Review) (int, error)
	// GetReview метод для получения отзыва по id.
	GetReview(context.Context, int) (models.Review, error)
	// GetReviewsByService метод для получения отзывов по id услуги.
	GetReviewsByService(context.Context, int) ([]models.Review, error)
	// GetReviewsByUser метод для получения отзывов по id пользователя.
	GetReviewsByUser(context.Context, int) ([]models.Review, error)
	// UpdateReview метод для обновления отзыва по id.
	// Метод также обновляет поле avg_score у услуги в случае изменения рейтинга отзыва.
	UpdateReview(context.Context, models.Review) error
	// BatchUpdateReviewsSentiment метод для батч-обновления полей Sentiment у отзывов.
	BatchUpdateReviewsSentiment(context.Context, []models.Review) error
	// DeleteReview метод для удаления отзыва по id.
	// Метод также обновляет поле avg_score у услуги в случае удаления отзыва.
	DeleteReview(context.Context, models.Review) error
	// GetUnsentimentedReviews возвращает необработанные отзывы без оценки настроения отзыва.
	GetUnsentimentedReviews(context.Context) ([]models.Review, error)
}
