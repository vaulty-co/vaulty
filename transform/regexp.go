package transform

import (
	"regexp"
	"sync"
)

type Regexp struct {
	Expression     string
	SubmatchNumber int
	Action         Transformer
	once           sync.Once
}

func (t *Regexp) Transform(body []byte) ([]byte, error) {
	var re *regexp.Regexp

	t.once.Do(func() {
		re = regexp.MustCompile(t.Expression)
	})

	// it does not make sence to do anything
	// if user specified submatch that does not exist
	if t.SubmatchNumber < 1 {
		return body, nil
	}

	result := re.FindSubmatchIndex(body)

	// if max position of submatch's end is
	// greater of max position of result it
	// means we don't have enough submatches
	if t.SubmatchNumber*2+1 > len(result)-1 {
		return body, nil
	}

	// result[2*n:2*n+1] identifies the indexes of the nth submatch
	n := t.SubmatchNumber
	prefix := body[0:result[2*n]]
	value := body[result[2*n]:result[2*n+1]]
	suffix := body[result[2*n+1]:]

	value, err := t.Action.Transform(value)
	if err != nil {
		return nil, err
	}

	newBody := append(prefix, value...)
	newBody = append(newBody, suffix...)

	return newBody, nil
}
