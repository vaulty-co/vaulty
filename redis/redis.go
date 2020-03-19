package redis

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func init() {
	client.WrapProcess(func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			before := time.Now()

			err := oldProcess(cmd)

			if err != nil {
				log.Println(`error running redis command "`, cmd, `": `, err)
			} else {
				log.Println(`running redis command "`, cmd, `" took `, time.Since(before))
			}

			return err
		}
	})
}

func Client() *redis.Client {
	return client
}
