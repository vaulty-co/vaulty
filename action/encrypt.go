package action

import "github.com/vaulty/vaulty/encryption"

type Encrypt struct {
	enc encryption.Encrypter
}

func (a *Encrypt) Transform(body []byte) ([]byte, error) {
	newBody, err := a.enc.Encrypt(body)
	if err != nil {
		return nil, err
	}
	return newBody, nil
}
