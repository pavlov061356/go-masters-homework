package postgres

import "pavlov061356/go-masters-homework/final_task/internal/models"

func (s *Storage) NewService(models.Service) (int, error) {
	return 0, nil
}

func (s *Storage) GetService(int) (models.Service, error) {
	return models.Service{}, nil
}

func (s *Storage) GetServices() ([]models.Service, error) {
	return nil, nil
}

func (s *Storage) UpdateService(models.Service) error {
	return nil
}

func (s *Storage) DeleteService(int) error {
	return nil
}
