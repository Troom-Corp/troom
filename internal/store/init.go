package store

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type InterfaceStore interface {
	Open()
	Close()
	Users() InterfaceUser
}

type store struct {
	db *sqlx.DB
}

func (s *store) Open() {
	database, _ := sqlx.Open("pgx", "postgres://linkedin_xhd9_user:IJlZf9iZnorMNMUsmtdfQqg0Vwkk3IHo@dpg-cl5r9ut6fh7c73etr82g-a.oregon-postgres.render.com/linkedin_xhd9")
	s.db = database
}

func (s *store) Close() {
	s.db.Close()
}

func (s *store) Users() InterfaceUser {
	return &user{db: s.db}
}

func NewStore() InterfaceStore {
	return &store{}
}
