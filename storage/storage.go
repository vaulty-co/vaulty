package storage

import "github.com/go-redis/redis"

type Storage struct {
	redisClient *redis.Client
}

func NewStorage(redisClient *redis.Client) *Storage {
	return &Storage{
		redisClient: redisClient,
	}
}
