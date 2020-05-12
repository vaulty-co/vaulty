package transform

type Transformer interface {
	Transform(body []byte) ([]byte, error)
}

type TransformerFunc func(body []byte) ([]byte, error)

func (f TransformerFunc) Transform(body []byte) ([]byte, error) {
	return f(body)
}
