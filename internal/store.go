package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// Store подключаемся к базе и возвращаем объект подключения
func Store() *pgx.Conn {
	var DatabaseUrl string = "postgres://postgres:123@localhost:5432/linkedin"
	conn, err := pgx.Connect(context.Background(), DatabaseUrl)
	if err != nil {
		panic(err)
	}
	return conn
}
