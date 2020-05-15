package action

import "github.com/rs/xid"

type Tokenize struct {
}

func (a *Tokenize) Transform(body []byte) ([]byte, error) {
	id, _ := xid.New().MarshalText()
	token := append([]byte("tkn"), id...)
	return token, nil
}
