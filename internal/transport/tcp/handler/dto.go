package handler

type messageResponse struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Difficulty int    `json:"difficulty,omitempty"`
}
