package storage

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var Sql Storage

type Storage struct {
	DB *sqlx.DB
}

func (s Storage) Open() (*sqlx.DB, error) {
	database, err := sqlx.Open("pgx", "postgres://postgres:123@localhost:5432/linkedin")
	if err != nil {
		return s.DB, err
	}
	s.DB = database
	return s.DB, nil
}

func (s Storage) Close() error {
	return s.DB.Close()
}
