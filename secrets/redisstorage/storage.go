package redisstorage

import (
	"github.com/go-redis/redis/v7"
	"github.com/vaulty/vaulty/encryption"
	"github.com/vaulty/vaulty/secrets"
)

var _ secrets.Storage = (*storage)(nil)

type storage struct {
	enc encryption.Encrypter
	db  *redis.Client
}

func Factory(conf *secrets.Config) (secrets.Storage, error) {
	return New(&Params{
		Encrypter: conf.Encrypter,
		RedisURL:  conf.StorageConfig.RedisURL,
	})
}

// Params is used as input to New.
type Params struct {
	// Enc is Encrypter that used to encrypt/decrypt data before putting it
	// into data storage
	Encrypter encryption.Encrypter

	// RedisURL is used for redis connection details. You can provide
	// the database, password and timeouts within the URL, e.g.
	// rediss://foo:bar@localhost:123
	RedisURL string
}

func New(params *Params) (secrets.Storage, error) {
	redisOptions, err := redis.ParseURL(params.RedisURL)
	if err != nil {
		return nil, err
	}

	db := redis.NewClient(redisOptions)

	// let's check if we can talk to redis
	if err := db.Ping().Err(); err != nil {
		return nil, err
	}

	return &storage{
		enc: params.Encrypter,
		db:  db,
	}, nil
}

// Set encrypts value and adds it to the map. It doesn't check for existing
// value and rewrites if it was set before.
func (s *storage) Set(key string, val []byte) error {
	encrypted, err := s.enc.Encrypt(val)
	if err != nil {
		return err
	}

	_, err = s.db.Set(key, encrypted, 0).Result()
	if err != nil {
		return err
	}

	return nil
}


// Set value and adds it to the map. It doesn't check for existing
// value and rewrites if it was set before.
func (s *storage) SetWithoutCrypto(key string, val string) error {
	s.db.Set(key, val, 0).Result()

	return nil
}

// Get decrypts the value from the map and returns it.
func (s *storage) Get(key string) ([]byte, error) {
	encrypted, err := s.db.Get(key).Result()
	if err != nil {
		return nil, err
	}

	decrypted, err := s.enc.Decrypt([]byte(encrypted))
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// Close closes the redis client, releasing any open resources
func (s *storage) Close() error {
	err := s.db.Close()
	return err
}
