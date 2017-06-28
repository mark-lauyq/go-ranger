package ranger_http

import (
	"encoding/json"
	"net/http"
)

// ResponseWriter struct
type ResponseWriter struct {
}

// ErrorResponse struct
type ErrorResponse struct {
	Status int `json:"status"`
	Data   struct {
		ErrorCode       string `json:"exception_type"`
		Message         string `json:"message"`
		Details         string `json:"developer_message"`
		MoreInformation string `json:"more_information"`
	} `json:"data"`
}

func (writer *ResponseWriter) writeErrorResponse(rw http.ResponseWriter, statusCode int, errorCode string, message string) {
	rw.WriteHeader(statusCode)

	json.NewEncoder(rw).Encode(ErrorResponse{
		Status: statusCode,
		Data: struct {
			ErrorCode       string `json:"exception_type"`
			Message         string `json:"message"`
			Details         string `json:"developer_message"`
			MoreInformation string `json:"more_information"`
		}{
			ErrorCode:       errorCode,
			Message:         message,
			Details:         "",
			MoreInformation: "null",
		},
	})
}
