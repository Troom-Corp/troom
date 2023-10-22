package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	New() *pgx.Conn
	Close(conn *pgx.Conn)
}

type DataBase struct {
}

func (d DataBase) New() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:123@localhost:5432/linkedin")
	if err != nil {
		panic(err)
	}
	return conn
}

func (d DataBase) Close(conn *pgx.Conn) {
	conn.Close(context.Background())
}

var SqlInterface Storage = DataBase{}
