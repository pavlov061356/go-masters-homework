package models

import "time"

// User - пользователь
type User struct {
	ID       int
	Username string
	Password string
}

// Service - предоставляемая услуга
type Service struct {
	ID          int
	Name        string
	Description string
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
