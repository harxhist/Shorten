package llm

import (
	"be/config"
	"be/constant"
	"be/helper"
	"be/logger"
	"be/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var log = logger.Logger



func SummarizeContent(content string) (string, time.Duration, error) {
    keyID, keyValue, err := config.LimitedConfig.GetRandomAPIKey("llm")
    if err != nil {
        return "", 0, fmt.Errorf("failed to get API Key: %v", err)
    }

    apiKey, ok := keyValue.(string)
    if !ok {
        return "", 0, fmt.Errorf("invalid API key format")
    }

    url := constant.APPCONFIG.GroqLLMEndpoint
    // url := constant.GROQ_LLM_ENDPOINT
    
    // Handle rate limiting
    resp,summaryLatency, err := makeRequest(url, content, apiKey)
    if err != nil {
        if err.Error() == "quota exceeded" {
            // Mark the key as exhausted and retry with a different key
            config.LimitedConfig.MarkKeyAsExhausted("llm", keyID)
            return SummarizeContent(content) // Recursive retry
        }
        return "",0, err
    }

    // Rest of your existing code...
    var groqResp model.GroqResponse
    if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
        return "",0,fmt.Errorf("resp decode error")
    }

    if len(groqResp.Choices) == 0 || groqResp.Choices[0].Message.Content == "" {
        return "",0, fmt.Errorf("empty response")
    }

    data := helper.ImproveSummary(groqResp.Choices[0].Message.Content)
    return data,summaryLatency, nil
}

func makeRequest(url, content, apiKey string) (*http.Response, time.Duration, error) {
    payload := map[string]interface{}{
        "messages": []map[string]interface{}{
            {
                "role":    constant.APPCONFIG.GroqUserRole,
                "content": constant.APPCONFIG.GroqLLMInstruction + content,
            },
        },
        "model": constant.APPCONFIG.GroqLLMModel,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
    }

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

    client := &http.Client{}
    start := time.Now()
    resp, err := client.Do(req)
    summaryLatency := helper.CalculateLatency(start)

    if err != nil {
        return nil, summaryLatency, fmt.Errorf("HTTP request failed: %w", err)
    }

    if resp.StatusCode == http.StatusTooManyRequests {
        return nil, summaryLatency, fmt.Errorf("quota exceeded")
    }

    if resp.StatusCode >= 400 && resp.StatusCode < 500 {
        return nil, summaryLatency, fmt.Errorf("client error: %s", resp.Status)
    }

    if resp.StatusCode >= 500 {
        return nil, summaryLatency, fmt.Errorf("server error: %s", resp.Status)
    }
    return resp, summaryLatency, nil
}
