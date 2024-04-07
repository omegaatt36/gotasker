package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// HTTPTestRequest defines the request for testing.
type HTTPTestRequest struct {
	// common
	Method string

	// server
	ServedURL   string
	HandleFuncs []gin.HandlerFunc

	// client
	RequestURLWithParams string
	Payload              map[string]interface{}
	Header               http.Header
}

// HTTPTestResponse defines the response for testing.
type HTTPTestResponse struct {
	StatusCode int
	Body       []byte
}

// PostForTest sends a POST request to the given URL with the given body and the given header. Put the route handler functions to last handleFuncs
func HTTPTest(req HTTPTestRequest) (*HTTPTestResponse, error) {
	var buffer io.Reader
	if req.Method == http.MethodPost ||
		req.Method == http.MethodPut ||
		req.Method == http.MethodPatch {
		bs, err := json.Marshal(req.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}

		buffer = bytes.NewBuffer(bs)
	}

	httpReq, err := http.NewRequest(req.Method, req.RequestURLWithParams, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, values := range req.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	httpReq.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()

	var serveFn func(string, ...gin.HandlerFunc) gin.IRoutes

	switch req.Method {
	case http.MethodGet:
		serveFn = router.GET
	case http.MethodPost:
		serveFn = router.POST
	case http.MethodPut:
		serveFn = router.PUT
	case http.MethodDelete:
		serveFn = router.DELETE
	default:
		return nil, fmt.Errorf("invalid method: %s", req.Method)
	}

	serveFn(req.ServedURL, req.HandleFuncs...)

	router.ServeHTTP(w, httpReq)

	return &HTTPTestResponse{
		StatusCode: w.Code,
		Body:       w.Body.Bytes(),
	}, nil
}
