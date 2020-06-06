package action

import "github.com/vaulty/vaulty/encrypt"

type Encrypt struct {
	enc encrypt.Encrypter
}

func (a *Encrypt) Transform(body []byte) ([]byte, error) {
	newBody, err := a.enc.Encrypt(body)
	if err != nil {
		return nil, err
	}
	return newBody, nil
}
