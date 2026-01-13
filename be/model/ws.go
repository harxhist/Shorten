package model

import "encoding/json"

type WebSocketRequest struct {
    Type      string          `json:"type"`       
    RequestID string          `json:"requestId"`  
    Payload   json.RawMessage `json:"payload"`    
}

type WebSocketResponseWithID struct {
	RequestID string      `json:"requestId"`
	Type      string      `json:"type"`
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}