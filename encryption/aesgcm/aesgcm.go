package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encryption"
)

var _ encryption.Encrypter = (*AesGcm)(nil)

func Factory(conf *config.Config) (encryption.Encrypter, error) {
	key := []byte(conf.Encryption.Key)

	// We use AES-256 which requires 32 bytes key
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: %d. Should be 32 bytes", len(key))
	}

	return NewEncrypter(key)
}

type AesGcm struct {
	block cipher.Block
}

func NewEncrypter(key []byte) (encryption.Encrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AesGcm{
		block: block,
	}, nil
}

// Encrypt returns hex encoded ciphertext
func (e *AesGcm) Encrypt(plaintext []byte) ([]byte, error) {
	aesgcm, err := cipher.NewGCM(e.block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	// pass nonce as dst to keep it within ciphertext
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	hexCiphertext := make([]byte, hex.EncodedLen(len(ciphertext)))
	hex.Encode(hexCiphertext, ciphertext)

	return hexCiphertext, nil
}

func (e *AesGcm) Decrypt(hexCiphertext []byte) ([]byte, error) {
	ciphertext := make([]byte, hex.DecodedLen(len(hexCiphertext)))
	_, err := hex.Decode(ciphertext, hexCiphertext)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(e.block)

	// ciphertext should include nonce
	if len(ciphertext) < aesgcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	nonce := ciphertext[:aesgcm.NonceSize()]
	ciphertext = ciphertext[aesgcm.NonceSize():]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)

	return plaintext, err
}
