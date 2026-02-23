package database

import (
	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	rdb *redis.Client
}

func NewRedisAdapter(addr string) *RedisAdapter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisAdapter{
		rdb: client,
	}
}
