package redisx

import "github.com/go-redis/redis"

type (
	Redis struct {
		client *redis.Client
	}
)
