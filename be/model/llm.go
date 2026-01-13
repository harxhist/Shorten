package model

// GroqResponse represents the structure of the API response from the summarization service.
type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}