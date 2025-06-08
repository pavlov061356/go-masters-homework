package postgres

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
	"time"
)

// NewService создает новую услугу и возвращает ее идентификатор.
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

// GetService возвращает услугу по ее идентификатору.
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

// GetServices возвращает все услуги.
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

// UpdateService обновляет информацию об услуге по ее идентификатору.
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

// DeleteService удаляет услугу по ее идентификатору.
func (s *Storage) DeleteService(ctx context.Context, id int) error {
	_, err := s.conn.ExecEx(ctx, `DELETE FROM services WHERE id = $1`, nil, id)

	if err != nil {
		return err
	}

	return nil
}

// RecomputeServicesScore пересчитывает средний рейтинг услуг.
func (s *Storage) RecomputeServicesScore(ctx context.Context) error {
	tx, err := s.conn.BeginEx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	_, err = tx.ExecEx(ctx, `
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

	_, err = tx.ExecEx(ctx,
		`UPDATE services_avg_score_compute_time SET last_avg_compute_time = DEFAULT;`,
		nil,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetLastRecomputeTime возвращает время последнего пересчета рейтингов услуг.
func (s *Storage) GetLastRecomputeTime(context.Context) (time.Time, error) {
	var time time.Time

	if err := s.conn.QueryRowEx(context.Background(),
		`SELECT last_avg_compute_time FROM services_avg_score_compute_time;`,
		nil,
	).Scan(
		&time,
	); err != nil {
		return time, err
	}

	return time, nil
}
