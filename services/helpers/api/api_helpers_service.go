package api_helpers_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	api_constants "go_boilerplate_project/constants/api"
	api_models "go_boilerplate_project/models/api"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ──────────────────────────────────────────────
// Request Parsing
// ──────────────────────────────────────────────

func (s *service) ParseJSONBody(c *gin.Context, dest any) []api_models.ValidationError {
	if err := c.ShouldBindJSON(dest); err != nil {
		return s.buildValidationErrors(err, dest, "json")
	}
	return nil
}

func (s *service) ParseQueryParams(c *gin.Context, dest any) []api_models.ValidationError {
	if err := c.ShouldBindQuery(dest); err != nil {
		return s.buildValidationErrors(err, dest, "form")
	}
	return nil
}

func (s *service) ParsePathParam(c *gin.Context, key string) string {
	return c.Param(key)
}

// ──────────────────────────────────────────────
// Response Helpers
// ──────────────────────────────────────────────

func (s *service) SendApiResponse(c *gin.Context, response *api_models.ApiResponse) {
	if response == nil {
		s.SendError(c, http.StatusInternalServerError, "internal server error", "response is nil")
		return
	}

	if response.StatusCode < 200 || response.StatusCode >= 600 {
		s.SendError(c, http.StatusInternalServerError, "internal server error", "invalid status code")
		return
	}

	if response.ValidationErrors != nil {
		s.SendValidationError(c, response.ValidationErrors)
		return
	}

	if response.Error != "" {
		s.SendError(c, response.StatusCode, response.Message, response.Error)
		return
	}

	if response.PageCount > 0 {
		s.SendPaginatedSuccess(c, response.StatusCode, response.Message, response.Data, &api_models.PaginationResult{
			Limit:      response.PageSize,
			Offset:     (response.PageNumber - 1) * response.PageSize,
			PageNumber: response.PageNumber,
			PageSize:   response.PageSize,
			TotalCount: response.TotalCount,
			TotalPages: response.TotalPages,
		})
		return
	}

	s.SendSuccess(c, response.StatusCode, response.Message, response.Data)
}

func (s *service) SendSuccess(c *gin.Context, statusCode int, message string, data any) {
	resp := api_models.ApiResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Timestamp:  time.Now(),
	}

	s.Input.Logger.Debugw("API response",
		"statusCode", statusCode,
		"message", message,
	)

	c.JSON(statusCode, resp)
}

func (s *service) SendError(c *gin.Context, statusCode int, message string, errDetail string) {
	resp := api_models.ApiResponse{
		StatusCode: statusCode,
		Message:    message,
		Error:      errDetail,
		Timestamp:  time.Now(),
	}

	s.Input.Logger.Warnw("API error response",
		"statusCode", statusCode,
		"message", message,
		"error", errDetail,
	)

	c.JSON(statusCode, resp)
}

func (s *service) SendValidationError(c *gin.Context, validationErrors []api_models.ValidationError) {
	resp := api_models.ApiResponse{
		StatusCode:       http.StatusBadRequest,
		Message:          "validation failed",
		ValidationErrors: validationErrors,
		Timestamp:        time.Now(),
	}

	s.Input.Logger.Warnw("API validation error response",
		"statusCode", http.StatusBadRequest,
		"errorCount", len(validationErrors),
		"errors", validationErrors,
	)

	c.JSON(http.StatusBadRequest, resp)
}

func (s *service) SendPaginatedSuccess(c *gin.Context, statusCode int, message string, data any, pagination *api_models.PaginationResult) {
	resp := api_models.ApiResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Timestamp:  time.Now(),
		PageCount:  countItems(data),
		PageNumber: pagination.PageNumber,
		PageSize:   pagination.PageSize,
		TotalCount: pagination.TotalCount,
		TotalPages: pagination.TotalPages,
	}

	s.Input.Logger.Debugw("API paginated response",
		"statusCode", statusCode,
		"message", message,
		"pageNumber", pagination.PageNumber,
		"pageSize", pagination.PageSize,
		"totalCount", pagination.TotalCount,
		"totalPages", pagination.TotalPages,
		"pageCount", resp.PageCount,
	)

	c.JSON(statusCode, resp)
}

// ──────────────────────────────────────────────
// Pagination
// ──────────────────────────────────────────────

func (s *service) GetPaginationFromQuery(c *gin.Context) (*api_models.PaginationInput, []api_models.ValidationError) {
	var errs []api_models.ValidationError

	pageNumber := api_constants.DefaultPageNumber
	pageSize := api_constants.DefaultPageSize

	if raw := c.Query("page_number"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_number",
				Message: "page_number must be a valid integer",
				Tag:     "numeric",
				Value:   raw,
			})
		} else if n < 1 {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_number",
				Message: "page_number must be at least 1",
				Tag:     "min",
				Value:   n,
			})
		} else {
			pageNumber = n
		}
	}

	if raw := c.Query("page_size"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_size",
				Message: "page_size must be a valid integer",
				Tag:     "numeric",
				Value:   raw,
			})
		} else if n < 1 {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_size",
				Message: "page_size must be at least 1",
				Tag:     "min",
				Value:   n,
			})
		} else if n > api_constants.MaxPageSize {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_size",
				Message: fmt.Sprintf("page_size must be at most %d", api_constants.MaxPageSize),
				Tag:     "max",
				Value:   n,
			})
		} else {
			pageSize = n
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &api_models.PaginationInput{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}, nil
}

func (s *service) CalculatePagination(input *api_models.PaginationInput, totalCount int) *api_models.PaginationResult {
	pageNumber := input.PageNumber
	pageSize := input.PageSize

	if pageNumber < 1 {
		pageNumber = api_constants.DefaultPageNumber
	}
	if pageSize < 1 {
		pageSize = api_constants.DefaultPageSize
	}
	if pageSize > api_constants.MaxPageSize {
		pageSize = api_constants.MaxPageSize
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + pageSize - 1) / pageSize
	}

	offset := (pageNumber - 1) * pageSize

	return &api_models.PaginationResult{
		Limit:      pageSize,
		Offset:     offset,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}

// ──────────────────────────────────────────────
// Validation Error Builder
// ──────────────────────────────────────────────

func (s *service) buildValidationErrors(err error, dest any, tagKey string) []api_models.ValidationError {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		result := make([]api_models.ValidationError, 0, len(validationErrs))
		for _, fe := range validationErrs {
			fieldName := resolveFieldName(dest, fe.StructField(), tagKey)
			result = append(result, api_models.ValidationError{
				Field:   fieldName,
				Message: buildFieldMessage(fieldName, fe),
				Tag:     fe.Tag(),
				Value:   fe.Value(),
			})
		}
		return result
	}

	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		field := unmarshalTypeErr.Field
		if field == "" {
			field = "body"
		}
		return []api_models.ValidationError{{
			Field:   field,
			Message: fmt.Sprintf("%s must be of type %s", field, unmarshalTypeErr.Type),
			Tag:     "type_mismatch",
			Value:   unmarshalTypeErr.Value,
		}}
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return []api_models.ValidationError{{
			Field:   "body",
			Message: fmt.Sprintf("malformed JSON: syntax error at position %d", syntaxErr.Offset),
			Tag:     "json_syntax",
		}}
	}

	if errors.Is(err, io.EOF) {
		return []api_models.ValidationError{{
			Field:   "body",
			Message: "request body is empty",
			Tag:     "required",
		}}
	}

	return []api_models.ValidationError{{
		Field:   "body",
		Message: err.Error(),
		Tag:     "parse_error",
	}}
}

// ──────────────────────────────────────────────
// Internal Helpers
// ──────────────────────────────────────────────

func resolveFieldName(obj any, structField string, tagKey string) string {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		if f, ok := t.FieldByName(structField); ok {
			if tag := f.Tag.Get(tagKey); tag != "" && tag != "-" {
				return strings.SplitN(tag, ",", 2)[0]
			}
		}
	}
	return structField
}

func buildFieldMessage(field string, fe validator.FieldError) string {
	param := fe.Param()
	isStr := fe.Kind() == reflect.String

	switch fe.Tag() {
	case "required", "required_if", "required_unless", "required_with", "required_without",
		"required_with_all", "required_without_all":
		return field + " is required"

	case "min":
		if isStr {
			return fmt.Sprintf("%s must be at least %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		if isStr {
			return fmt.Sprintf("%s must be at most %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "len":
		if isStr {
			return fmt.Sprintf("%s must be exactly %s characters long", field, param)
		}
		return fmt.Sprintf("%s must have exactly %s items", field, param)

	case "eq":
		return fmt.Sprintf("%s must be equal to %s", field, param)
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, param)

	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, param)
	case "nefield":
		return fmt.Sprintf("%s must not match %s", field, param)
	case "gtfield":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "ltfield":
		return fmt.Sprintf("%s must be less than %s", field, param)

	case "email":
		return field + " must be a valid email address"
	case "url":
		return field + " must be a valid URL"
	case "uri":
		return field + " must be a valid URI"
	case "uuid", "uuid3", "uuid4", "uuid5":
		return field + " must be a valid UUID"
	case "ip":
		return field + " must be a valid IP address"
	case "ipv4":
		return field + " must be a valid IPv4 address"
	case "ipv6":
		return field + " must be a valid IPv6 address"
	case "json":
		return field + " must be valid JSON"
	case "jwt":
		return field + " must be a valid JWT token"
	case "e164":
		return field + " must be a valid E.164 phone number"

	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "alphanumunicode":
		return field + " must contain only unicode letters and numbers"
	case "numeric":
		return field + " must be numeric"
	case "boolean":
		return field + " must be a boolean"
	case "lowercase":
		return field + " must be lowercase"
	case "uppercase":
		return field + " must be uppercase"

	case "contains":
		return fmt.Sprintf("%s must contain '%s'", field, param)
	case "excludes":
		return fmt.Sprintf("%s must not contain '%s'", field, param)
	case "startswith":
		return fmt.Sprintf("%s must start with '%s'", field, param)
	case "endswith":
		return fmt.Sprintf("%s must end with '%s'", field, param)

	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime in format '%s'", field, param)

	default:
		return fmt.Sprintf("%s failed on '%s' validation", field, fe.Tag())
	}
}

func countItems(data any) int {
	if data == nil {
		return 0
	}
	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return v.Len()
	}
	return 0
}
