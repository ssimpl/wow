package book

const (
	messageTypeQuote     = "quote"
	messageTypeChallenge = "challenge"
	messageTypeError     = "error"
)

type messageResponse struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Difficulty int    `json:"difficulty,omitempty"`
}
