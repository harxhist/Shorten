package model

// Application Request Format
type Request struct {
	URL string `json:"url"`
}
// Application Response Format
type Response struct {
	Status      string       `json:"status"`
	Content     interface{}  `json:"responseBody"`
	ErrorDetail *ErrorDetail `json:"errorDetail"`
}
// Application Error Format
type ErrorDetail struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
// Application Feeback Format
type Feedback struct {
	RequestID         string  `json:"requestID"`
	Feedback          int     `json:"feedback"`
	FeedbackText      string  `json:"feedbackText"`
	TTSPlayedDuration float64 `json:"ttsplayedDuration"`
	NavigatedToClean  bool    `json:"navigatedToClean"`
}
