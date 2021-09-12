package redisx

import "github.com/go-redis/redis"
//
//var RedisSrv
//
//func init(){
//	redisSrv := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "123456",
//		DB:       0,
//	})
//	RedisSrv = redisSrv
//}

type (
	RedisI interface {

	}
)

func NewRedisService(opt *redis.Options) *Redis{
	client := redis.NewClient(opt)
	return &Redis{
		client: client,
	}
}

func (s *Redis) Subcribe(){

}
