package model

type Config struct {
	LokiURL                   string   `json:"LOKI_URL"`
	AppName                   string   `json:"APP_NAME"`
	Environment               string   `json:"ENVIRONMENT"`
	AllowedOrigin             []string `json:"ALLOWED_ORIGIN"`
	DBHost                    string   `json:"DB_HOST"`
	DBName                    string   `json:"DB_NAME"`
	DBUser                    string   `json:"DB_USER"`
	DBPass                    string   `json:"DB_PASS"`
	DBSecretName              string   `json:"DB_SECRET_NAME"`
	S3Endpoint                string   `json:"S3_ENDPOINT"`
	S3Region                  string   `json:"S3_REGION"`
	S3ID                      string   `json:"S3_ID"`
	S3Secret                  string   `json:"S3_SECRET"`
	S3Token                   string   `json:"S3_TOKEN"`
	S3AudioBucket             string   `json:"S3_AUDIO_BUCKET"`
	S3SpeechMarkBucket        string   `json:"S3_SPEECH_MARK_BUCKET"`
	S3LLMBucket               string   `json:"S3_LLM_BUCKET"`
	S3CleanedBucket           string   `json:"S3_CLEANED_BUCKET"`
	GroqLLMEndpoint           string   `json:"GROQ_LLM_ENDPOINT"`
	GroqUserRole              string   `json:"GROQ_USER_ROLE"`
	GroqLLMModel              string   `json:"GROQ_LLM_MODEL"`
	GroqLLMInstruction        string   `json:"GROQ_LLM_INSTRUCTION"`
	AWSRegion                 string   `json:"AWS_REGION"`
	AWSVoiceID                string   `json:"AWS_VOICE_ID"`
	AWSSpeechMarkOutputFormat string   `json:"AWS_SPEECH_MARK_OUTPUT_FORMAT"`
	AWSSpeechMarkType         string   `json:"AWS_SPEECH_MARK_TYPE"`
	AWSAudioOutputFormat      string   `json:"AWS_AUDIO_OUTPUT_FORMAT"`
	PublicKey                 string   `json:"PUBLIC_KEY"`
}