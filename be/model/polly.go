package model

type PollyRequest struct {
	ReadText string `json:"readText"`
}

type PollyResponse struct {
	AudioData       string          `json:"audioData"`
	SpeechMarksData []SpeechMarkData `json:"speechMarksData"`
}

type SpeechMarkData struct {
	Time     int64  `json:"time"`
	Type     string `json:"type"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Value    string `json:"value"`
}