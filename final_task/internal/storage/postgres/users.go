package postgres

import "pavlov061356/go-masters-homework/final_task/internal/models"

func (s *Storage) NewUser(user models.User) (int, error) {
	return 0, nil
}

func (s *Storage) GetUser(id int) (models.User, error) {
	return models.User{}, nil
}

func (s *Storage) GetUserByEmail(email string) (models.User, error) {
	return models.User{}, nil
}

func (s *Storage) UpdateUser(user models.User) error {
	return nil
}

func (s *Storage) DeleteUser(id int) error {
	return nil
}
