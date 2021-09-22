package redisx

import (
	"github.com/go-redis/redis/v8"
	"orchestrator/internal/transactioncache"
)

type (
	Service struct {
		redisClient *redis.Client
		txCacheSrv transactioncache.ServiceI
	}
)