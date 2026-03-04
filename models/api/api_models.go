package api_models

import "time"

type ApiResponse struct {
	StatusCode       int               `json:"status_code"`
	Message          string            `json:"message"`
	MessageKey       string            `json:"-"` // Translation key; when set, Message is translated from this
	Error            string            `json:"error,omitempty"`
	ErrorKey         string            `json:"-"` // Translation key for Error; when set, Error is translated
	ErrorKeyParams   map[string]string  `json:"-"` // Params for ErrorKey template (e.g. {"Detail": "..."})
	Data             any               `json:"data,omitempty"`
	Timestamp        time.Time         `json:"timestamp"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
	PageCount        int               `json:"page_count,omitempty"`
	PageNumber       int               `json:"page_number,omitempty"`
	PageSize         int               `json:"page_size,omitempty"`
	TotalCount       int               `json:"total_count,omitempty"`
	TotalPages       int               `json:"total_pages,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Value   any    `json:"value,omitempty"`
}

type PaginationInput struct {
	PageNumber int
	PageSize   int
}

type PaginationResult struct {
	Limit      int
	Offset     int
	PageNumber int
	PageSize   int
	TotalCount int
	TotalPages int
}
