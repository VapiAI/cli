package voice

import (
	"time"
)

// APIHandler manages API request/response logging and handling
type APIHandler struct {
	client      *VoiceClient
	requestLog  chan APIRequest
	responseLog chan APIResponse
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(client *VoiceClient) *APIHandler {
	return &APIHandler{
		client:      client,
		requestLog:  make(chan APIRequest, 100),
		responseLog: make(chan APIResponse, 100),
	}
}

// LogRequest logs an API request
func (h *APIHandler) LogRequest(method, url string, headers map[string]string, body interface{}) {
	req := APIRequest{
		Method:    method,
		URL:       url,
		Headers:   headers,
		Body:      body,
		Timestamp: time.Now(),
	}

	select {
	case h.requestLog <- req:
	default:
		// Channel full, skip logging
	}
}

// LogResponse logs an API response
func (h *APIHandler) LogResponse(statusCode int, headers map[string]string, body interface{}, duration time.Duration) {
	resp := APIResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
		Duration:   duration,
		Timestamp:  time.Now(),
	}

	select {
	case h.responseLog <- resp:
	default:
		// Channel full, skip logging
	}
}

// GetRequestLog returns the request log channel
func (h *APIHandler) GetRequestLog() <-chan APIRequest {
	return h.requestLog
}

// GetResponseLog returns the response log channel
func (h *APIHandler) GetResponseLog() <-chan APIResponse {
	return h.responseLog
}

// FormatRequest formats an API request for display
func FormatRequest(req APIRequest) string {
	return req.Timestamp.Format("15:04:05") + " → " + req.Method + " " + req.URL
}

// FormatResponse formats an API response for display
func FormatResponse(resp APIResponse) string {
	return resp.Timestamp.Format("15:04:05") + " ← " + 
		   string(rune(resp.StatusCode)) + " " + 
		   resp.Duration.String()
}