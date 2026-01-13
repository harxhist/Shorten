package response

import (
	"encoding/json"
	"net/http"
	"be/model"
)

// SendErrorResponse sends an error response with the specified error code and message.
func SendErrorResponse(w http.ResponseWriter, statusCode int, errorCode, errorMessage string) {
	response := model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SendSuccessResponse sends a successful response with the provided data.
func SendSuccessResponse(w http.ResponseWriter, responseData interface{}) {
	response := model.Response{
		Status:       "SUCCESS",
		Content:   responseData,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func WebSocketErrorResponseWithID(requestID, message string) model.WebSocketResponseWithID {
	return model.WebSocketResponseWithID{
		RequestID: requestID,
		Status:    "error",
		Message:   message,
	}
}

func WebSocketSuccessResponseWithID(requestID, responseType string, data interface{}) model.WebSocketResponseWithID {
	return model.WebSocketResponseWithID{
		RequestID: requestID,
		Type:      responseType,
		Status:    "success",
		Data:      data,
	}
}
