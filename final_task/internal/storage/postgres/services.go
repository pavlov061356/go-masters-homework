package postgres

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
)

func (s *Storage) NewService(ctx context.Context, service models.Service) (int, error) {
	var id int
	err := s.conn.QueryRowEx(ctx,
		`INSERT INTO services (name, description) VALUES ($1, $2) RETURNING id`,
		nil,
		service.Name,
		service.Description,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetService(ctx context.Context, id int) (models.Service, error) {
	var service models.Service

	err := s.conn.QueryRowEx(ctx,
		`SELECT name, description, avg_score FROM services WHERE id = $1`,
		nil,
		id,
	).Scan(
		&service.Name,
		&service.Description,
		&service.AvgScore,
	)
	if err != nil {
		return models.Service{}, err
	}

	service.ID = id

	return service, nil
}

func (s *Storage) GetServices(context.Context) ([]models.Service, error) {
	var services []models.Service

	rows, err := s.conn.QueryEx(context.Background(),
		`SELECT id, name, description, avg_score FROM services`,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var service models.Service

		if err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.AvgScore,
		); err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}

func (s *Storage) UpdateService(ctx context.Context, service models.Service) error {
	_, err := s.conn.ExecEx(ctx,
		`UPDATE services SET name = $1, description = $2 WHERE id = $3`,
		nil,
		service.Name,
		service.Description,
		service.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteService(ctx context.Context, id int) error {
	_, err := s.conn.ExecEx(ctx, `DELETE FROM services WHERE id = $1`, nil, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RecomputeServicesScore(ctx context.Context) error {
	_, err := s.conn.ExecEx(ctx, `
	UPDATE services 
		SET avg_score = 
			(SELECT COALESCE(AVG(r.score), 0)
			FROM reviews r
			WHERE r.service_id = services.id
			)
	`,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}
