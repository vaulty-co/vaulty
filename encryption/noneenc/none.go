package noneenc

import (
	"bytes"
	"encoding/base64"

	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encryption"
)

var _ encryption.Encrypter = (*None)(nil)

func Factory(conf *config.Config) (encryption.Encrypter, error) {
	return &None{}, nil
}

// None encryptor is base64 encoder. It does not encrypt any data
// Should be used only for demo / dev purposes.
type None struct {
}

const demoNotice string = "(demo encryption)"

func New() encryption.Encrypter {
	return new(None)
}

func (*None) Encrypt(plaintext []byte) ([]byte, error) {
	encoded := make([]byte, base64.RawStdEncoding.EncodedLen(len(plaintext)))

	base64.RawStdEncoding.Encode(encoded, plaintext)

	return append(encoded, []byte(demoNotice)...), nil
}

func (*None) Decrypt(ciphertext []byte) ([]byte, error) {
	ciphertext = bytes.TrimSuffix(ciphertext, []byte(demoNotice))

	decoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(ciphertext)))
	_, err := base64.RawStdEncoding.Decode(decoded, ciphertext)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
