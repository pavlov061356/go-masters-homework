package postgres

import "pavlov061356/go-masters-homework/final_task/internal/models"

func (s *Storage) NewReview(models.Review) (int, error) {
	return 0, nil
}
func (s *Storage) GetReview(int) (models.Review, error) {
	return models.Review{}, nil
}
func (s *Storage) GetReviewsByService(int) ([]models.Review, error) {
	return nil, nil
}
func (s *Storage) GetReviewsByUser(int) ([]models.Review, error) {
	return nil, nil
}
func (s *Storage) UpdateReview(models.Review) error {
	return nil
}
func (s *Storage) DeleteReview(int) error {
	return nil
}
