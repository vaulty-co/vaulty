package core

import "github.com/go-redis/redis"

func NewRedisClient() *redis.Client {
	redisOptions, err := redis.ParseURL(Config.Redis.URL)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(redisOptions)
}
