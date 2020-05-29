package secrets

type SecretsStorage interface {
	Set(key string, val []byte) error
	Get(key string) ([]byte, error)
}
