package postgres

import (
	"context"
	"errors"
	"pavlov061356/go-masters-homework/final_task/internal/consts"
	"pavlov061356/go-masters-homework/final_task/internal/models"

	"github.com/jackc/pgx"
)

// NewReview метод для создания отзыва.
// Обновляет среднюю оценку услуги и возвращает id созданного отзыва.
func (s *Storage) NewReview(ctx context.Context, review models.Review) (int, error) {
	tx, err := s.conn.BeginEx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.RollbackEx(ctx)

	var id int

	err = tx.QueryRowEx(ctx,
		`INSERT INTO reviews (
			content,
			score,
			service_id,
			reviewer_id
		) VALUES ($1, $2, $3, $4)
		RETURNING id;`,
		nil,
		review.Content,
		review.Score,
		review.ServiceID,
		review.ReviewerID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecEx(ctx, `
		UPDATE services s
			SET avg_score = (
			SELECT 
				CASE 
					WHEN COUNT(r.id) > 1 
					THEN (COALESCE(s.avg_score, 0) * (COUNT(r.id) - 1) + $2) / (COUNT(r.id))
					ELSE $2
				END
			FROM reviews r
			WHERE r.service_id = s.id
			)
		WHERE s.id = $1
		`,
		nil,
		review.ServiceID,
		review.Score,
	)
	if err != nil {
		return 0, err
	}

	if err := tx.CommitEx(ctx); err != nil {
		return 0, err
	}

	return id, nil
}

// GetReview метод для получения отзыва по id.
func (s *Storage) GetReview(ctx context.Context, id int) (models.Review, error) {
	var review models.Review
	err := s.conn.QueryRowEx(ctx,
		`SELECT
			content,
			sentiment,
			score,
			created_at,
			reviewer_id,
			service_id
		FROM reviews
		WHERE id = $1
		`,
		nil,
		id,
	).Scan(
		&review.Content,
		&review.Sentiment,
		&review.Score,
		&review.CreatedAt,
		&review.ReviewerID,
		&review.ServiceID,
	)

	if err != nil {
		return models.Review{}, err
	}

	review.ID = id

	return review, nil
}

// GetReviewsByService метод для получения отзывов по id услуги.
func (s *Storage) GetReviewsByService(ctx context.Context, id int) ([]models.Review, error) {
	rows, err := s.conn.QueryEx(ctx,
		`SELECT
			id,
			content,
			sentiment,
			score,
			created_at,
			reviewer_id,
			service_id
		FROM reviews
		WHERE service_id = $1
		ORDER BY created_at DESC
		`,
		nil,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var review models.Review

		err := rows.Scan(
			&review.ID,
			&review.Content,
			&review.Sentiment,
			&review.Score,
			&review.CreatedAt,
			&review.ReviewerID,
			&review.ServiceID,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// GetReviewsByUser метод для получения отзывов по id пользователя.
func (s *Storage) GetReviewsByUser(ctx context.Context, id int) ([]models.Review, error) {
	rows, err := s.conn.QueryEx(ctx,
		`SELECT
			id,
			content,
			sentiment,
			score,
			created_at,
			reviewer_id,
			service_id
		FROM reviews
		WHERE reviewer_id = $1
		ORDER BY created_at DESC
		`,
		nil,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var review models.Review

		err := rows.Scan(
			&review.ID,
			&review.Content,
			&review.Sentiment,
			&review.Score,
			&review.CreatedAt,
			&review.ReviewerID,
			&review.ServiceID,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// UpdateReview метод для обновления отзыва по id.
// Метод также обновляет поле avg_score у услуги в случае изменения рейтинга отзыва.
func (s *Storage) UpdateReview(ctx context.Context, review models.Review) error {
	tx, err := s.conn.BeginEx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.RollbackEx(ctx)

	var currentScore int
	err = tx.QueryRowEx(ctx,
		`SELECT score FROM reviews WHERE id = $1`,
		nil,
		review.ID,
	).Scan(&currentScore)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	_, err = tx.ExecEx(ctx,
		`UPDATE reviews SET  
			content = $1,
			score = $2,
			service_id = $3,
			reviewer_id = $4
		WHERE id = $5
		`,
		nil,
		review.Content,
		review.Score,
		review.ServiceID,
		review.ReviewerID,
		review.ID,
	)
	if err != nil {
		return err
	}

	if currentScore != review.Score {
		_, err = tx.ExecEx(ctx,
			`WITH count AS (
				SELECT COUNT(*) AS count
				FROM reviews
    			WHERE service_id = $1
			) 
			UPDATE services
			SET avg_score = (
			    (avg_score * (SELECT count FROM count) - $2 + $3) / 
			    (SELECT count FROM count)
			)
			WHERE id = $1
		`,
			nil,
			review.ServiceID,
			currentScore,
			review.Score,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.CommitEx(ctx); err != nil {
		return err
	}

	return nil
}

// BatchUpdateReviewsSentiment метод для батч-обновления полей Sentiment у отзывов.
func (s *Storage) BatchUpdateReviewsSentiment(ctx context.Context, reviews []models.Review) error {
	batch := s.conn.BeginBatch()

	for _, review := range reviews {
		batch.Queue(
			`UPDATE reviews SET sentiment = $1 WHERE id = $2`,
			[]any{
				review.Sentiment,
				review.ID,
			},
			nil,
			nil,
		)
	}

	if err := batch.Send(ctx, nil); err != nil {
		return err
	}

	return nil
}

// DeleteReview метод для удаления отзыва по id.
// Метод также обновляет поле avg_score у услуги в случае удаления отзыва.
func (s *Storage) DeleteReview(ctx context.Context, review models.Review) error {
	tx, err := s.conn.BeginEx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.RollbackEx(ctx)

	_, err = tx.ExecEx(ctx,
		`DELETE FROM reviews WHERE id = $1`,
		nil,
		review.ID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecEx(ctx, `
		UPDATE services s
			SET avg_score = (
    			SELECT 
    			    CASE 
    			        WHEN COUNT(r.id) > 0 
    			        THEN (avg_score * (COUNT(r.id) + 1) - $2) / COUNT(r.id)
    			        ELSE 0
    			    END
    			FROM reviews r
    			WHERE r.service_id = s.id
			)
		WHERE s.id = $1
		`,
		nil,
		review.ServiceID,
		review.Score,
	)
	if err != nil {
		return err
	}

	if err := tx.CommitEx(ctx); err != nil {
		return err
	}

	return nil
}

// GetUnsentimentedReviews возвращает необработанные отзывы без оценки настроения отзыва.
func (s *Storage) GetUnsentimentedReviews(ctx context.Context) ([]models.Review, error) {
	rows, err := s.conn.QueryEx(ctx, `
	SELECT id, content
	FROM reviews
	WHERE sentiment = $1
	`,
		nil,
		consts.SentimentUnknown,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		err = rows.Scan(&r.ID, &r.Content)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, r)
	}

	return reviews, nil
}
