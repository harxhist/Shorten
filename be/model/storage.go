package model

import (
	"time"
)

type DBRow struct {
	RequestID         string    //B (Primary Key)
	RequestTime       time.Time //B
	RequestedURL      string    //B
	IPAddress         string    //B
	AudioLink         string    //B
	LLMResponseLink   string    //B
	SpeechMarksLink   string    //B
	CleanedTextLink   string    //B
	NumberOfWords     int       //B
	TotalLatency      float64   //B
	SummaryLatency    float64   //B
	TTSLatency        float64   //B
	Feedback          int       //F
	FeedbackText      string    //F
	TTSPlayedDuration float64   //F
	NavigatedToClean  bool      //F
}


type StorageItem struct {
    Name     string
    Store    func() error
}