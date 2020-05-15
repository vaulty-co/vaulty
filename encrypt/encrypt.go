package encrypt

import "fmt"

type Encrypter interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

func NewEncrypter(key string) (Encrypter, error) {
	if key == "" {
		return &None{}, nil
	}

	// We use AES-256 which requires 32 bytes key
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: %d. Should be 32 bytes", len(key))
	}

	return NewAesGcm(key)
}
