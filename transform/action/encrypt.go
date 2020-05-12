package action

type Encrypt struct {
}

func (a *Encrypt) Transform(body []byte) ([]byte, error) {
	newBody := append(body, " encrypted"...)
	return newBody, nil
}
