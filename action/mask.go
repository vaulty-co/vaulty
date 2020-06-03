package action

import (
	"bytes"
)

type Mask struct {
	Symbol []byte
}

func (a *Mask) Transform(body []byte) ([]byte, error) {
	replacer := a.Symbol
	if a.Symbol == nil {
		replacer = []byte("*")
	}
	newBody := bytes.Repeat(replacer, len(body))
	return newBody, nil
}
