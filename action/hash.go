package action

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hash struct {
	salt string
}

func (a *Hash) Transform(body []byte) ([]byte, error) {
	// add salt to body
	body = append(body, []byte(a.salt)...)

	sum := sha256.Sum256(body)
	newBody := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(newBody, sum[:])

	return newBody, nil
}
