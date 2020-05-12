package action

import "github.com/rs/xid"

type Tokenize struct {
}

func (action *Tokenize) Transform(body []byte) ([]byte, error) {
	id, _ := xid.New().MarshalText()
	token := append([]byte("tkn"), id...)
	return token, nil
}
