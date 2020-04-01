package transformer

// Serializable structure for Sidekiq Worker
// result with raw body
type Result struct {
	Body []byte `json:"body"`
}
