package task

// Serializable structure for Sidekiq Worker
// with raw body
type Response struct {
	Body []byte `json:"body"`
}
