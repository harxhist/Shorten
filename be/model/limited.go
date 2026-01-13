package model

type LimitedConfig struct {
	LLM           LLM               `json:"llm"`
	TTS           TTS               `json:"tts"`
}

type LLM struct {
	Keys          map[string]string `json:"keys"`
	ExhaustedKeys map[string]bool   `json:"-"`
}

type TTS struct {
	Keys          map[string]TTSKey `json:"keys"`
	ExhaustedKeys map[string]bool   `json:"-"`
}

type TTSKey struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}