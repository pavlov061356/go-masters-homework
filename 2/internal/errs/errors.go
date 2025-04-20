package errs

import (
	"fmt"
	"math/rand"
)

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "данные не найдены"
}

type ErrBadRequest struct{}

func (e *ErrBadRequest) Error() string {
	return "неверный запрос"
}

type ErrConflict struct {
	name string
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("сущность с именем %s уже существует", e.name)
}

type ErrInternalServerError struct {
	message string
}

func (e *ErrInternalServerError) Error() string {
	return fmt.Sprintf("внутренняя ошибка сервера. Сообщение: %s", e.message)
}

// RandomErr возвращает случайную ошибку
// Вспомогательная функция для примера
func RandomErr() error {
	number := rand.Intn(4)

	switch number {
	case 0:
		return &ErrNotFound{}
	case 1:
		return &ErrBadRequest{}
	case 2:
		return &ErrConflict{name: "test"}
	case 3:
		return &ErrInternalServerError{message: "test"}
	default:
		return &ErrInternalServerError{message: "default"}
	}
}
