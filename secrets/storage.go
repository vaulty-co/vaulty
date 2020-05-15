package secrets

type SecretStorage interface {
	Set(key string, val []byte) error
	Get(key string) ([]byte, error)
}
