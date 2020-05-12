package transform

type Transformer interface {
	Transform(body []byte) ([]byte, error)
}
