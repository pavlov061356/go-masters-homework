package models

import (
	"errors"
	"time"
)

// User - пользователь
type User struct {
	ID       int
	Email    string
	Username string
	Password string
}

// Service - предоставляемая услуга
type Service struct {
	ID          int
	Name        string
	Description string
	AvgScore    float64
}

func (s Service) Validate() error {
	if s.Name == "" {
		return errors.New("название услуги не может быть пустым")
	}
	if s.Description == "" {
		return errors.New("описание услуги не может быть пустым")
	}
	return nil
}

// Review - отзыв о предоставляемой услуге
type Review struct {
	ID         int
	Content    string
	Sentiment  int // 0 - не определён, 1 - положительный, 2 - нормальный, 3 - отрицательный
	Score      int // 1-5 оценка пользователя
	CreatedAt  time.Time
	ReviewerID int
	ServiceID  int
}

func (r Review) Validate() error {
	if r.Content == "" {
		return errors.New("отзыв не может быть пустым")
	}

	if r.Score < 1 || r.Score > 5 {
		return errors.New("оценка должна быть в диапазоне от 1 до 5 включительно")
	}

	if r.ReviewerID == 0 {
		return errors.New("отзыв должен содержать идентификатор пользователя")
	}

	if r.ServiceID == 0 {
		return errors.New("отзыв должен содержать идентификатор услуги")
	}
	return nil
}
