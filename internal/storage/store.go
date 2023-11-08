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
	database, err := sqlx.Open("pgx", "postgres://linkedin_xhd9_user:IJlZf9iZnorMNMUsmtdfQqg0Vwkk3IHo@dpg-cl5r9ut6fh7c73etr82g-a.oregon-postgres.render.com/linkedin_xhd9")
	if err != nil {
		return s.DB, err
	}
	s.DB = database
	return s.DB, nil
}

func (s Storage) Close() error {
	return s.DB.Close()
}
