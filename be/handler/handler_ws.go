package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"context"
	"github.com/gorilla/websocket"
	"be/response"
	"be/model"
	"be/helper"
	"be/llm"
	"be/storage"
	"be/constant"
	"strings"
)


var upgrader = websocket.Upgrader{
	// Restrict origins to configured allowed origins
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		for _, allowed := range constant.APPCONFIG.AllowedOrigin {
			if strings.EqualFold(origin, allowed) {
				return true
			}
		}
		return false
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to upgrade to WebSocket: ", err)
		response.SendErrorResponse(w, http.StatusBadRequest, "WS-400", "Failed to upgrade to WebSocket")
		return
	}

	for {

		var msg model.WebSocketRequest
		if err := conn.ReadJSON(&msg); err != nil {
			log.Error("Invalid WS request: ", err)
			conn.WriteJSON(response.WebSocketErrorResponseWithID(msg.RequestID, "Invalid request format"))
			break
		}
		
		if msg.Type == "process" {
			
			go handleSummarizeAndTTS(conn, msg.RequestID, msg.Payload)
		} else {
			log.Error("Unsupported request type: ", msg.Type)
			conn.WriteJSON(response.WebSocketErrorResponseWithID(msg.RequestID, "Unsupported request type"))
		}
	}
}

func handleSummarizeAndTTS(conn *websocket.Conn, requestID string, payload json.RawMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	start := time.Now()

	var summarizeReq model.Request
	if err := json.Unmarshal(payload, &summarizeReq); err != nil {
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Invalid request payload"))
		return
	}

	parseStart := time.Now()
	parsedURL, err := helper.ParseURL(summarizeReq.URL)
	parseLatency := time.Since(parseStart)
	if err != nil {
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to parse URL: "+err.Error()))
		return
	}

	fetchStart := time.Now()
	content, err := fetchContent(parsedURL)
	fetchLatency := time.Since(fetchStart)
	if err != nil {
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to fetch content: "+err.Error()))
		return
	}

	cleanStart := time.Now()
	tagsToRemove := []string{"script", "iframe", "noscript", "style"}
	content, err = helper.CleanContent(content, tagsToRemove)
	cleanLatency := time.Since(cleanStart)
	if err != nil {
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to clean content: "+err.Error()))
		return
	}

	extractStart := time.Now()
	content, err = helper.ExtractText(content)
	extractLatency := time.Since(extractStart)
	if err != nil {
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to extract text: "+err.Error()))
		return
	}

	summarizeStart := time.Now()
	content, _, err = llm.SummarizeContent(content)
	summaryLatency := time.Since(summarizeStart)
	if err != nil {
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to summarize content: "+err.Error()))
		return
	}
	LLMResp := content

	if err := conn.WriteJSON(response.WebSocketSuccessResponseWithID(requestID, "summarize", content)); err != nil {
		log.Error("Failed to send summarize response: ", err)
		return
	}

	markdownExtractStart := time.Now()
	cleanedMarkdown := helper.ExtractTextFromMarkdown(content)
	markdownExtractLatency := time.Since(markdownExtractStart)

	ttsStart := time.Now()
	ttsReq := &model.PollyRequest{ReadText: string(cleanedMarkdown)}
	audioData, _, err := handlePollyRequest(ctx, ttsReq)
	ttsLatency := time.Since(ttsStart)
	if err != nil {
		log.Error("Failed to generate audio: ", err)
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to generate audio: "+err.Error()))
		return
	}

	if err := conn.WriteJSON(response.WebSocketSuccessResponseWithID(requestID, "tts", audioData)); err != nil {
		log.Error("Failed to send TTS response: ", err)
		return
	}

	totalLatency := time.Since(start)
	
	log.Info("Total Latency: ", totalLatency, 
		" Parse URL Latency: ", parseLatency, 
		" Fetch Content Latency: ", fetchLatency, 
		" Clean Content Latency: ", cleanLatency, 
		" Extract Text Latency: ", extractLatency, 
		" Summarization Latency: ", summaryLatency, 
		" Markdown Extraction Latency: ", markdownExtractLatency, 
		" TTS Latency: ", ttsLatency)

	// Store Data
	if err := storage.StoreData(requestID, audioData.AudioData, LLMResp, cleanedMarkdown, audioData.SpeechMarksData); err != nil {
		log.Error("Failed to store data: ", err)
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to store data"))
		return
	}

	// Generate Presigned URLs
	mp3Link, err := storage.GeneratePresignedURL("audio", requestID, 604800*time.Second)
	if err != nil {
		log.Error("Failed to generate audio presigned URL: ", err)
	}

	llmResponseLink, err := storage.GeneratePresignedURL("llm", requestID, 604800*time.Second)
	if err != nil {
		log.Error("Failed to generate LLM response presigned URL: ", err)
	}

	speechMarksLink, err := storage.GeneratePresignedURL("spm", requestID, 604800*time.Second)
	if err != nil {
		log.Error("Failed to generate speech marks presigned URL: ", err)
	}

	cleanedTextLink, err := storage.GeneratePresignedURL("text", requestID, 604800*time.Second)
	if err != nil {
		log.Error("Failed to generate cleaned text presigned URL: ", err)
	}

	// Save request metadata to database
	dbRow := &model.DBRow{
		RequestID:       requestID,
		RequestedURL:    summarizeReq.URL,
		IPAddress:       conn.RemoteAddr().String(),
		RequestTime:     start,
		AudioLink:       mp3Link,
		LLMResponseLink: llmResponseLink,
		SpeechMarksLink: speechMarksLink,
		CleanedTextLink: cleanedTextLink,
		NumberOfWords:   len(audioData.SpeechMarksData),
		TotalLatency:    float64(totalLatency.Milliseconds()),
		SummaryLatency:  float64(summaryLatency.Milliseconds()),
		TTSLatency:      float64(ttsLatency.Milliseconds()),
	}

	if err := storage.InsertRequest(ctx, dbRow); err != nil {
		log.Error("Failed to store request data in database: ", err)
		conn.WriteJSON(response.WebSocketErrorResponseWithID(requestID, "Failed to store request metadata"))
		return
	}

	conn.Close()
}