package network_models

import (
	"net/http"
	"time"
)

// HTTPMethod represents the HTTP request method.
type HTTPMethod string

const (
	HTTPMethodGet     HTTPMethod = "GET"
	HTTPMethodPost    HTTPMethod = "POST"
	HTTPMethodPut     HTTPMethod = "PUT"
	HTTPMethodPatch   HTTPMethod = "PATCH"
	HTTPMethodDelete  HTTPMethod = "DELETE"
	HTTPMethodHead    HTTPMethod = "HEAD"
	HTTPMethodOptions HTTPMethod = "OPTIONS"
)

// FetchInput contains all parameters for an HTTP fetch request.
type FetchInput struct {
	// Route is the full URL or path to call (required).
	Route string

	// Method is the HTTP method (required).
	Method HTTPMethod

	// Headers are optional request headers. Key is header name, value is header value.
	Headers map[string]string

	// Payload is the raw request body as any type. Use for POST, PUT, PATCH requests.
	Payload any

	// QueryParams are optional URL query parameters. Key is param name, value is param value.
	QueryParams map[string]string

	// Timeout is the request timeout. If zero, a default is used.
	Timeout time.Duration

	// ResponseModel is an optional pointer to unmarshal the response body into (e.g. &myStruct).
	// If provided and response is JSON, the parsed model will be in FetchOutput.ParsedModel.
	ResponseModel any

	// SkipTLSVerify when true skips TLS certificate verification (use for dev only).
	SkipTLSVerify bool
}

// FetchOutput contains the complete HTTP response details.
type FetchOutput struct {
	// StatusCode is the HTTP status code (e.g. 200, 404, 500).
	StatusCode int

	// Status is the HTTP status string (e.g. "200 OK").
	Status string

	// Headers contains all response headers. Key is header name (canonical), value is slice of values.
	Headers http.Header

	// BodyBytes is the raw response body as bytes.
	BodyBytes []byte

	// ParsedModel is the response body unmarshaled into the type pointed to by FetchInput.ResponseModel.
	// Only populated when FetchInput.ResponseModel was provided and unmarshaling succeeded.
	ParsedModel any

	// Duration is the total time the request took.
	Duration time.Duration

	// RequestURL is the final URL that was called (including query params).
	RequestURL string

	// ContentLength is the Content-Length from response headers, or -1 if unknown.
	ContentLength int64
}
