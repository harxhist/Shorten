package handler

import (
        "context"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "io"
        "strings"
        "sync"
        "time"
        "bytes"
        "be/config"
        "be/constant"
        "be/helper"
        "be/model"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/credentials"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/polly"
)

// Result structs for goroutines
type audioResult struct {
        audioSrc string
        err      error
}

type speechMarkResult struct {
        marks []model.SpeechMarkData
        err   error
}

var (
        pollyClient     *polly.Polly
        pollyClientMutex sync.Mutex
        sessionCache    sync.Map
        sessionTTL      = 15 * time.Minute
)

func handlePollyRequest(ctx context.Context, req *model.PollyRequest) (*model.PollyResponse, time.Duration, error) {
        start := time.Now()

        // Initialize Polly client (thread-safe)
        err := initializePollyClient()
        if err != nil {
                return nil, helper.CalculateLatency(start), err
        }

        audioChan := make(chan audioResult, 1)
        speechMarkChan := make(chan speechMarkResult, 1)

        go generateAudio(ctx, pollyClient, req.ReadText, audioChan)
        go generateSpeechMarks(ctx, pollyClient, req.ReadText, speechMarkChan)

        var audioRes audioResult
        var speechMarkRes speechMarkResult

        for i := 0; i < 2; i++ {
                select {
                case <-ctx.Done():
                        return nil, helper.CalculateLatency(start), ctx.Err()
                case audioRes = <-audioChan:
                case speechMarkRes = <-speechMarkChan:
                }
        }

        if audioRes.err != nil {
                return nil, helper.CalculateLatency(start), fmt.Errorf("audio generation failed: %w", audioRes.err)
        }

        if speechMarkRes.err != nil {
                return nil, helper.CalculateLatency(start), fmt.Errorf("speech marks generation failed: %w", speechMarkRes.err)
        }

        pollyResp := model.PollyResponse{
                AudioData:       audioRes.audioSrc,
                SpeechMarksData: speechMarkRes.marks,
        }

        return &pollyResp, helper.CalculateLatency(start), nil
}

func initializePollyClient() error {
        if pollyClient != nil {
                return nil
        }

        pollyClientMutex.Lock()
        defer pollyClientMutex.Unlock()

        if pollyClient != nil {
                return nil
        }

        keyID, keyValue, err := config.LimitedConfig.GetRandomAPIKey("tts")
        if err != nil {
                return fmt.Errorf("failed to get AWS credentials: %w", err)
        }

        ttsKey, ok := keyValue.(map[string]interface{})
        if !ok {
                return fmt.Errorf("invalid TTS key format")
        }

        accessKey, _ := ttsKey["accessKey"].(string)
        secretKey, _ := ttsKey["secretKey"].(string)

        sess, err := createAWSSession(accessKey, secretKey)
        if err != nil {
                config.LimitedConfig.MarkKeyAsExhausted("tts", keyID)
                return fmt.Errorf("failed to create AWS session: %w", err)
        }

        pollyClient = polly.New(sess)
        return nil
}

func generateAudio(ctx context.Context, client *polly.Polly, text string, resultChan chan<- audioResult) {
        output, err := client.SynthesizeSpeechWithContext(ctx, &polly.SynthesizeSpeechInput{
                OutputFormat: aws.String(constant.APPCONFIG.AWSAudioOutputFormat),
                Text:        aws.String(text),
                VoiceId:     aws.String(constant.APPCONFIG.AWSVoiceID),
        })

        if err != nil {
                resultChan <- audioResult{err: err}
                return
        }
        defer output.AudioStream.Close()

        var buf bytes.Buffer
        if _, err := io.Copy(&buf, output.AudioStream); err != nil {
                resultChan <- audioResult{err: err}
                return
        }

        resultChan <- audioResult{audioSrc: fmt.Sprintf("data:audio/mp3;base64,%s", base64.StdEncoding.EncodeToString(buf.Bytes()))}
}

func generateSpeechMarks(ctx context.Context, client *polly.Polly, text string, resultChan chan<- speechMarkResult) {
        output, err := client.SynthesizeSpeechWithContext(ctx, &polly.SynthesizeSpeechInput{
                OutputFormat:    aws.String(constant.APPCONFIG.AWSSpeechMarkOutputFormat),
                Text:           aws.String(text),
                VoiceId:        aws.String(constant.APPCONFIG.AWSVoiceID),
                SpeechMarkTypes: []*string{aws.String(constant.APPCONFIG.AWSSpeechMarkType)},
        })

        if err != nil {
                resultChan <- speechMarkResult{err: err}
                return
        }
        defer output.AudioStream.Close()

        var marks []model.SpeechMarkData
        decoder := json.NewDecoder(output.AudioStream)
        for decoder.More() {
                var mark model.SpeechMarkData
                if err := decoder.Decode(&mark); err != nil {
                        if err != io.EOF {
                                resultChan <- speechMarkResult{err: err}
                                return
                        }
                        break
                }
                marks = append(marks, mark)
        }

        resultChan <- speechMarkResult{marks: marks}
}

func createAWSSession(accessKey, secretKey string) (*session.Session, error) {
        cacheKey := fmt.Sprintf("%s:%s", accessKey, secretKey)

        if cachedSession, ok := sessionCache.Load(cacheKey); ok {
                return cachedSession.(*session.Session), nil
        }

        sess, err := session.NewSession(&aws.Config{
                Region:      aws.String(constant.APPCONFIG.AWSRegion),
                Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        })

        if err != nil {
                return nil, fmt.Errorf("failed to create AWS session: %w", err)
        }

        sessionCache.Store(cacheKey, sess)

        go func() {
                time.Sleep(sessionTTL)
                sessionCache.Delete(cacheKey)
        }()

        return sess, nil
}

func isQuotaExceeded(err error) bool {
        errStr := strings.ToLower(err.Error())
        return strings.Contains(errStr, "quota exceeded") ||
                strings.Contains(errStr, "throttlingexception") ||
                strings.Contains(errStr, "rate exceeded")
}