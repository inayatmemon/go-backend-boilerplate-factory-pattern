package network_service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	network_models "go_boilerplate_project/models/network"
)

const defaultTimeout = 30 * time.Second

// Fetch executes an HTTP request and returns the full response details.
func (s *service) Fetch(ctx context.Context, input *network_models.FetchInput) (*network_models.FetchOutput, error) {
	if input == nil {
		s.Input.Logger.Error("Fetch input is nil")
		return nil, fmt.Errorf("fetch input is nil")
	}

	start := time.Now()
	s.Input.Logger.Infow("Fetch request started",
		"route", input.Route,
		"method", input.Method,
	)

	// Build URL with query params
	requestURL, err := s.buildURL(input.Route, input.QueryParams)
	if err != nil {
		s.Input.Logger.Errorw("Failed to build request URL",
			"route", input.Route,
			"error", err,
		)
		return nil, fmt.Errorf("build url: %w", err)
	}

	// Create request body reader
	var bodyReader io.Reader
	if input.Payload != nil {
		payload, ok := input.Payload.([]byte)
		if !ok {
			paylodJSON, err := json.Marshal(input.Payload)
			if err != nil {
				s.Input.Logger.Errorw("Failed to marshal payload",
					"payload", input.Payload,
					"error", err,
				)
				return nil, fmt.Errorf("marshal payload: %w", err)
			}
			payload = paylodJSON
			ok = true
		}
		bodyReader = bytes.NewReader(payload)
		s.Input.Logger.Debugw("Request payload attached",
			"size", len(payload),
		)
	}

	req, err := http.NewRequestWithContext(ctx, string(input.Method), requestURL, bodyReader)
	if err != nil {
		s.Input.Logger.Errorw("Failed to create HTTP request",
			"url", requestURL,
			"error", err,
		)
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	for k, v := range input.Headers {
		req.Header.Set(k, v)
	}
	if len(input.Headers) > 0 {
		s.Input.Logger.Debugw("Request headers set", "count", len(input.Headers))
	}

	// Configure client
	timeout := input.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: input.SkipTLSVerify,
			},
		},
	}

	s.Input.Logger.Debugw("Executing HTTP request", "url", requestURL)
	resp, err := client.Do(req)
	if err != nil {
		s.Input.Logger.Errorw("HTTP request failed",
			"url", requestURL,
			"error", err,
		)
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Input.Logger.Errorw("Failed to read response body",
			"statusCode", resp.StatusCode,
			"error", err,
		)
		return nil, fmt.Errorf("read body: %w", err)
	}

	duration := time.Since(start)
	output := &network_models.FetchOutput{
		StatusCode:    resp.StatusCode,
		Status:        resp.Status,
		Headers:       resp.Header.Clone(),
		BodyBytes:     bodyBytes,
		Duration:      duration,
		RequestURL:    requestURL,
		ContentLength: resp.ContentLength,
	}

	// Parse into ResponseModel if provided
	if input.ResponseModel != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, input.ResponseModel); err != nil {
			s.Input.Logger.Warnw("Failed to unmarshal response into target model",
				"statusCode", resp.StatusCode,
				"bodySize", len(bodyBytes),
				"error", err,
			)
			// Don't fail - we still have BodyBytes and raw response
		} else {
			output.ParsedModel = input.ResponseModel
			s.Input.Logger.Debugw("Response unmarshaled successfully",
				"statusCode", resp.StatusCode,
			)
		}
	}

	s.Input.Logger.Infow("Fetch request completed",
		"url", requestURL,
		"statusCode", resp.StatusCode,
		"duration", duration,
		"bodySize", len(bodyBytes),
	)

	return output, nil
}

// buildURL appends query parameters to the base route.
func (s *service) buildURL(route string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return route, nil
	}

	parsed, err := url.Parse(route)
	if err != nil {
		return "", err
	}

	q := parsed.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}
