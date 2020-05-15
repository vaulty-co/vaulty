package secrets

import "github.com/vaulty/proxy/encrypt"

type ephemeral struct {
	enc  encrypt.Encrypter
	data map[string][]byte
}

func NewEphemeralStorage(enc encrypt.Encrypter) SecretStorage {
	return &ephemeral{
		data: make(map[string][]byte),
		enc:  enc,
	}
}

func (s *ephemeral) Set(key string, val []byte) error {
	encrypted, err := s.enc.Encrypt(val)
	if err != nil {
		return err
	}

	s.data[key] = encrypted

	return nil
}

func (s *ephemeral) Get(key string) ([]byte, error) {
	encrypted := s.data[key]
	decryted, err := s.enc.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	return decryted, nil
}
