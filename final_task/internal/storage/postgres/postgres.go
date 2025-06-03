package postgres

import (
	"pavlov061356/go-masters-homework/final_task/internal/storage"

	"github.com/jackc/pgx"
)

var _ storage.Interface = &Storage{}

type Storage struct {
	conn *pgx.Conn
}

func New(dsn string) (*Storage, error) {
	config, err := pgx.ParseConnectionString(dsn)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(config)
	if err != nil {
		return nil, err
	}
	return &Storage{
		conn: conn,
	}, nil
}
