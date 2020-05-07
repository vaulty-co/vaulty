package storage

import "github.com/go-redis/redis"

type redisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) Storage {
	return &redisStorage{
		redisClient: redisClient,
	}
}

func (s *redisStorage) Reset() {
	s.redisClient.FlushAll()
}
