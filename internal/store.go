package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// StoreSetup подключаемся к базе и возвращаем объект подключения
func StoreSetup() *pgx.Conn {
	var DatabaseUrl string = "postgres://postgres:123@localhost:5432/linkedin"
	conn, err := pgx.Connect(context.Background(), DatabaseUrl)
	if err != nil {
		panic(err)
	}
	return conn
}

var Store = StoreSetup()
