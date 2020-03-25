package core

import (
	"github.com/go-redis/redis"
)

func NewRedisClient(config *Configuration) *redis.Client {
	redisOptions, err := redis.ParseURL(config.Redis.URL)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(redisOptions)
}
