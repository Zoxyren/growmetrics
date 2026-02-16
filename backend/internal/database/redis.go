package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(addr string) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return Rdb.Ping(Ctx).Err()
}
