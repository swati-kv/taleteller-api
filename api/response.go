package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
	RequestID string      `json:"request_id,omitempty"`

	// can be used when returning multiple form errors
	Errors []ErrorInfo `json:"errors,omitempty"`
}

// ErrorInfo specifies what info are we sending.
// Use IsEmpty method instead of comparing with struct literal.
type ErrorInfo struct {
	Field    string          `json:"field"`
	Message  string          `json:"message"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// IsEmpty returns if the field and message is empty or not. Does not check the metadata as it is dynamic
func (e ErrorInfo) IsEmpty() bool {
	return e.Field == "" && e.Message == ""
}

func FormattedErrors(errors []ErrorInfo) (message string) {
	for _, err := range errors {
		message += fmt.Sprintf("%s %s ", err.Field, err.Message)
	}

	return
}

func RespondWithJSON(rw http.ResponseWriter, status int, response Response) {
	respBytes, err := json.Marshal(response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(respBytes)
}

func RespondWithRawJSON(rw http.ResponseWriter, status int, response interface{}) {
	respBytes, err := json.Marshal(response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(respBytes)
}

func RespondWithError(rw http.ResponseWriter, status int, response Response) {
	requestID := rw.Header().Get("X-Request-ID")
	if requestID == "" {
		requestID = "not-set"
	}

	response.RequestID = requestID

	respBytes, err := json.Marshal(response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(respBytes)
}
