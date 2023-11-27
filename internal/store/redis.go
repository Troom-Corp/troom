package store

import (
	"github.com/redis/go-redis/v9"
)

var Redis RedisDB

type RedisDB struct {
	DB *redis.Client
}

func (r RedisDB) Open() *redis.Client {
	r.DB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return r.DB
}
