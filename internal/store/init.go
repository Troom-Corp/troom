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
	database, _ := sqlx.Open("pgx", "postgres://troomdb_user:D4Lj9R3yT03HrpjqqGzj7oZcUDMVD64z@dpg-clk8cg58td7s73dd5vfg-a.frankfurt-postgres.render.com/troomdb")
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
