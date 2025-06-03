-- Сервисные таблицы

CREATE TABLE migrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Общие таблицы

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE service (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE review (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL DEFAULT '', -- Текст отзыва
    sentiment INTEGER DEFAULT 0 CHECK (sentiment >= 0 AND sentiment <= 3), -- 0 - не определён, 1 - положительный, 2 - нормальный, 3 - отрицательный 
    score INTEGER NOT NULL CHECK (score > 0 AND score <= 5), -- Численная оценка, оставленная пользователем
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    reviewer_id INTEGER REFERENCES users(id),
    service_id INTEGER REFERENCES service(id)
);

CREATE INDEX idx_review_reviewer_created_at ON review (reviewer_id, created_at); -- b-tree индекс
