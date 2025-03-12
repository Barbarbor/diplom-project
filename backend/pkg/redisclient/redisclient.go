package redisclient

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func Init(address, password string, db int) {
	Client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
}
