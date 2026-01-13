package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"be/model"
)


func PutInternalServerErrorResponse(errMsg interface{}, rw http.ResponseWriter) {
	// defer metrics.IncToInternalErrorCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    "SH1-500",
			ErrorMessage: fmt.Sprintf("Internal Server Error, %v", errMsg),
		},
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutBadRequestErrorResponse(errMsg interface{}, rw http.ResponseWriter) {
	// defer metrics.IncToBadRequestErrorCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    "SH1-400",
			ErrorMessage: fmt.Sprintf("Invalid request, %v", errMsg),
		},
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutInvalidJwtResponse(errMsg interface{}, rw http.ResponseWriter) {
	// defer metrics.IncToInvalidJwtCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    "SH1-401",
			ErrorMessage: fmt.Sprintf("Invalid Authorization token, %v", errMsg),
		},
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutNotFoundResponse(errMsg interface{}, rw http.ResponseWriter) {
	// defer metrics.IncToInvalidJwtCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    "SH1-404",
			ErrorMessage: fmt.Sprintf("Invalid URL, %v", errMsg),
		},
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutQuotaExceededResponse(errMsg interface{}, rw http.ResponseWriter) {
	// defer metrics.IncToInvalidJwtCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    "SH1-429",
			ErrorMessage: fmt.Sprintf("Quota Exceeded: try after some time, %v", errMsg),
		},
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutUpstreamInputErrorResponse(errMsg interface{}, rw http.ResponseWriter) {
	// defer metrics.IncToUnprocessableErrorCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    "SH1-422",
			ErrorMessage: fmt.Sprintf("Upstream Input Invalid or unprocessable entity, %v", errMsg),
		},
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutUpstreamOutputErrorResponse(errMsg interface{}, rw http.ResponseWriter) {
	// defer metrics.IncNoContentErrorCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status: "FAILURE",
		ErrorDetail: &model.ErrorDetail{
			ErrorCode:    "SH1-503",
			ErrorMessage: fmt.Sprintf("Upstream Output Error, %v", errMsg),
		},
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutSuccessResponse(respData interface{}, rw http.ResponseWriter) {
	// defer metrics.IncToRequestSuccessCounter()
	respBytes, _ := json.Marshal(model.Response{
		Status:       "SUCCESS",
		Content: 	respData,
	})
	PutHandlerResponse(rw, io.NopCloser(bytes.NewReader(respBytes)), http.StatusOK)
}

func PutHandlerResponse(rw http.ResponseWriter, respReader io.ReadCloser, statusCode int) {
	defer respReader.Close()
	header := make(map[string][]string)
	// header[constant.CONTENT_TYPE] = []string{constant.JSON_HEADER}
	for key, vals := range header {
		for _, val := range vals {
			rw.Header().Add(key, val)
		}
	}
	rw.WriteHeader(statusCode)
	io.Copy(rw, respReader)
}
