package redis

import "github.com/go-redis/redis"

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func Client() *redis.Client {
	return client
}
