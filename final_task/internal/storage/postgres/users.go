package postgres

import (
	"context"
	"pavlov061356/go-masters-homework/final_task/internal/models"
)

// NewUser создаёт нового пользователя и возвращает его идентификатор.
func (s *Storage) NewUser(ctx context.Context, user models.User) (int, error) {
	var id int

	err := s.conn.QueryRowEx(ctx,
		`INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`,
		nil,
		user.Email,
		user.Password,
		user.Username,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser возвращает пользователя по идентификатору.
func (s *Storage) GetUser(ctx context.Context, id int) (models.User, error) {
	var user models.User
	err := s.conn.QueryRowEx(ctx,
		`SELECT email, name, password FROM users WHERE id = $1 LIMIT 1
		`,
		nil,
		id,
	).Scan(
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		return user, err
	}

	user.ID = id

	return user, nil
}

// GetUserByEmail возвращает пользователя по email.
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := s.conn.QueryRowEx(ctx,
		`SELECT id, name, password FROM users WHERE email = $1 LIMIT 1
		`,
		nil,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		return user, err
	}

	user.Email = email

	return user, nil
}

// UpdateUser обновляет данные пользователя.
func (s *Storage) UpdateUser(ctx context.Context, user models.User) error {
	_, err := s.conn.ExecEx(ctx,
		`UPDATE users SET email = $1, name = $2, password = $3 WHERE id = $4`,
		nil,
		user.Email,
		user.Username,
		user.Password,
		user.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

// DeleteUser удаляет пользователя по идентификатору.
func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	_, err := s.conn.ExecEx(ctx, `DELETE FROM users WHERE id = $1`, nil, id)
	if err != nil {
		return err
	}
	return nil
}
