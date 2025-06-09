package models

import (
	"errors"
	"time"
)

var (
	ErrEmptyServiceName = errors.New("название услуги не может быть пустым")
	ErrEmptyServiceDesc = errors.New("описание услуги не может быть пустым")
)

var (
	ErrEmptyReviewName  = errors.New("отзыв не может быть пустым")
	ErrEmptyReviewScore = errors.New("оценка должна быть в диапазоне от 1 до 5 включительно")
	ErrEmptyReviewerID  = errors.New("отзыв должен содержать идентификатор пользователя")
	ErrEmptyServiceID   = errors.New("отзыв должен содержать идентификатор услуги")
)

// User - пользователь.
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Service - предоставляемая услуга.
type Service struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	AvgScore    float64 `json:"avg_score"`
}

func (s Service) Validate() error {
	if s.Name == "" {
		return ErrEmptyServiceName
	}

	if s.Description == "" {
		return ErrEmptyServiceDesc
	}

	return nil
}

// Review - отзыв о предоставляемой услуге.
type Review struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	Sentiment  int       `json:"sentiment"` // 0 - не определён, 1 - положительный, 2 - нормальный, 3 - отрицательный.
	Score      int       `json:"score"`     // 1-5 оценка пользователя.
	CreatedAt  time.Time `json:"created_at"`
	ReviewerID int       `json:"reviewer_id"`
	ServiceID  int       `json:"service_id"`
}

func (r Review) Validate() error {
	if r.Content == "" {
		return ErrEmptyReviewName
	}

	if r.Score < 1 || r.Score > 5 {
		return ErrEmptyReviewScore
	}

	if r.ReviewerID == 0 {
		return ErrEmptyReviewerID
	}

	if r.ServiceID == 0 {
		return ErrEmptyServiceID
	}

	return nil
}
