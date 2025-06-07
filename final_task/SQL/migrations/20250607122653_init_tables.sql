-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE INDEX users_email ON users (email);

INSERT INTO users (id, email, name, password) VALUES (1, 'test@test.test', 'test', 'testPassword');

ALTER SEQUENCE users_id_seq RESTART WITH 100;

CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    avg_score REAL DEFAULT 0,
    last_avg_compute_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO services (id, name, description) VALUES (1, 'test_service', 'test_service description');

ALTER SEQUENCE services_id_seq RESTART WITH 100;

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL DEFAULT '', -- Текст отзыва
    sentiment INTEGER DEFAULT 0 CHECK (sentiment >= 0 AND sentiment <= 3), -- 0 - не определён, 1 - положительный, 2 - нормальный, 3 - отрицательный 
    score INTEGER NOT NULL CHECK (score > 0 AND score <= 5), -- Численная оценка, оставленная пользователем
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    reviewer_id INTEGER REFERENCES users(id),
    service_id INTEGER REFERENCES services(id)
);

CREATE INDEX idx_review_reviewer_created_at ON reviews (reviewer_id, created_at); -- b-tree индекс

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS services;

DROP TABLE IF EXISTS reviews;
-- +goose StatementEnd
