package transformer

import (
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/vaulty/proxy/core"
)

var redisClient *redis.Client

func TestMain(m *testing.M) {
	core.LoadConfig("../config/test.yml")

	redisOptions, err := redis.ParseURL(core.Config.Redis.URL)
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(redisOptions)
	redisClient.FlushAll()
	exitCode := m.Run()
	redisClient.FlushAll()
	os.Exit(exitCode)
}
