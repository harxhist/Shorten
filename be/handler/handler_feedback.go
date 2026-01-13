package handler

import (
    "be/model"
    "be/response"
	"be/storage"
    "encoding/json"
    "net/http"
	"strings"
	"reflect"
	"fmt"
)

func HandleFeedback(w http.ResponseWriter, r *http.Request) {
    // Check if method is PUT
    if r.Method != http.MethodPut {
        response.PutBadRequestErrorResponse("Only PUT method is allowed", w)
        return
    }

    // Decode the request body
    var feedback model.Feedback
    if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
        response.PutBadRequestErrorResponse("Invalid request payload", w)
        return
    }

    // Validate the feedback
    if err := validateFeedback(feedback); err != nil {
        response.PutBadRequestErrorResponse(err.Error(), w)
        return
    }

    // Update feedback in database
    if err := updateFeedbackInDB(feedback); err != nil {
        handleFeedbackError(err, w)
        return
    }

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    response.PutSuccessResponse("Feedback updated successfully", w)
}

func validateFeedback(feedback model.Feedback) error {
    
	if(reflect.TypeOf(feedback.RequestID).Kind() != reflect.String){
		return fmt.Errorf("requestID must be a string")
	}
	if(reflect.TypeOf(feedback.Feedback).Kind() != reflect.Int){
		return fmt.Errorf("feedback must be an integer")
	}	
	if(reflect.TypeOf(feedback.FeedbackText).Kind() != reflect.String){
		return fmt.Errorf("feedbackText must be a string")
	}
	if(reflect.TypeOf(feedback.TTSPlayedDuration).Kind() != reflect.Float64){
		return fmt.Errorf("ttsplayedDuration must be a float64")
	}
	if(reflect.TypeOf(feedback.NavigatedToClean).Kind() != reflect.Bool){
		return fmt.Errorf("navigatedToClean must be a boolean")
	}
	if feedback.RequestID == "" {
        return fmt.Errorf("requestID is required")
    }
    if feedback.Feedback < -1 || feedback.Feedback > 1 {
        return fmt.Errorf("feedback must be -1, 0, or 1")
    }
	
    return nil
}

func updateFeedbackInDB(feedback model.Feedback) error {
	err := storage.UpdateFeedback(feedback)
	if(err != nil){
		log.Error("Failed to update feedback in database: ", err)
	}
    return nil;
}

func handleFeedbackError(err error, w http.ResponseWriter) {
    switch {
    case strings.Contains(err.Error(), "no record found with RequestID"):
        log.Error("No record found with given requestID")
        response.PutNotFoundResponse("No record found with given requestID", w)
    case err.Error() == "duplicate key value violates unique constraint":
        log.Error("Feedback already exists for this request")
        response.PutBadRequestErrorResponse("Feedback already exists for this request", w)
    default:
        log.Error("Failed to update feedback: ", err)
        response.PutInternalServerErrorResponse("Failed to update feedback", w)
    }
}