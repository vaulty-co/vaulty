package action

import "github.com/vaulty/proxy/encrypt"

type Decrypt struct {
	enc encrypt.Encrypter
}

func (a *Decrypt) Transform(body []byte) ([]byte, error) {
	newBody, err := a.enc.Decrypt(body)
	if err != nil {
		return nil, err
	}
	return newBody, nil
}
