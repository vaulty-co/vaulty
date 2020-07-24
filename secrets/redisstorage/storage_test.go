package redisstorage

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/encryption/noneenc"
)

var db *redis.Client
var redisURL string

func TestMain(m *testing.M) {
	var err error

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "3.2", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		redisURL = fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp"))

		db = redis.NewClient(&redis.Options{
			Addr: redisURL,
		})

		return db.Ping().Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// When you're done, kill and remove the container
	err = pool.Purge(resource)

	os.Exit(code)
}

func TestRedisStorage(t *testing.T) {
	encrypter := noneenc.New()

	st, err := New(&Params{
		RedisURL:  "redis://" + redisURL,
		Encrypter: encrypter,
	})
	require.NoError(t, err)

	err = st.Set("key", []byte("value"))
	require.NoError(t, err)

	// let's check that inside db we store encrypted value
	rawVal, err := db.Get("key").Result()
	require.NoError(t, err)
	// "demo encryption" is added by noneenc encryptor
	require.Contains(t, rawVal, "demo encryption")

	val, err := st.Get("key")
	require.Equal(t, []byte("value"), val)

	err = st.Close()
	require.NoError(t, err)
}
