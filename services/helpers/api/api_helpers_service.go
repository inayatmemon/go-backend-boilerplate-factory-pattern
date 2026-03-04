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
		lang := s.getLanguage(c)
		return s.buildValidationErrors(err, dest, "json", lang)
	}
	return nil
}

func (s *service) ParseQueryParams(c *gin.Context, dest any) []api_models.ValidationError {
	if err := c.ShouldBindQuery(dest); err != nil {
		lang := s.getLanguage(c)
		return s.buildValidationErrors(err, dest, "form", lang)
	}
	return nil
}

func (s *service) getLanguage(c *gin.Context) string {
	if s.Input.Services.Multilang != nil {
		return s.Input.Services.Multilang.GetLanguage(c)
	}
	return "en"
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

	lang := s.getLanguage(c)
	message := s.resolveMessage(lang, response)
	errorDetail := s.resolveError(lang, response)

	if response.ValidationErrors != nil {
		s.SendValidationError(c, response.ValidationErrors)
		return
	}

	if errorDetail != "" {
		s.SendError(c, response.StatusCode, message, errorDetail)
		return
	}

	if response.PageCount > 0 {
		s.SendPaginatedSuccess(c, response.StatusCode, message, response.Data, &api_models.PaginationResult{
			Limit:      response.PageSize,
			Offset:     (response.PageNumber - 1) * response.PageSize,
			PageNumber: response.PageNumber,
			PageSize:   response.PageSize,
			TotalCount: response.TotalCount,
			TotalPages: response.TotalPages,
		})
		return
	}

	s.SendSuccess(c, response.StatusCode, message, response.Data)
}

func (s *service) resolveMessage(lang string, r *api_models.ApiResponse) string {
	if r.MessageKey != "" && s.Input.Services.Multilang != nil {
		return s.Input.Services.Multilang.GetMessage(lang, r.MessageKey, nil)
	}
	return r.Message
}

func (s *service) resolveError(lang string, r *api_models.ApiResponse) string {
	if r.ErrorKey != "" && s.Input.Services.Multilang != nil {
		return s.Input.Services.Multilang.GetMessage(lang, r.ErrorKey, r.ErrorKeyParams)
	}
	return r.Error
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
	lang := s.getLanguage(c)
	msg := s.translate(lang, "validation_failed", nil)
	resp := api_models.ApiResponse{
		StatusCode:       http.StatusBadRequest,
		Message:          msg,
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
	lang := s.getLanguage(c)

	pageNumber := api_constants.DefaultPageNumber
	pageSize := api_constants.DefaultPageSize

	if raw := c.Query("page_number"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_number",
				Message: s.translate(lang, "validation_page_number_int", nil),
				Tag:     "numeric",
				Value:   raw,
			})
		} else if n < 1 {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_number",
				Message: s.translate(lang, "validation_page_number_min", nil),
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
				Message: s.translate(lang, "validation_page_size_int", nil),
				Tag:     "numeric",
				Value:   raw,
			})
		} else if n < 1 {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_size",
				Message: s.translate(lang, "validation_page_size_min", nil),
				Tag:     "min",
				Value:   n,
			})
		} else if n > api_constants.MaxPageSize {
			errs = append(errs, api_models.ValidationError{
				Field:   "page_size",
				Message: s.translate(lang, "validation_page_size_max", map[string]string{"param": fmt.Sprintf("%d", api_constants.MaxPageSize)}),
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

func (s *service) buildValidationErrors(err error, dest any, tagKey, lang string) []api_models.ValidationError {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		result := make([]api_models.ValidationError, 0, len(validationErrs))
		for _, fe := range validationErrs {
			fieldName := resolveFieldName(dest, fe.StructField(), tagKey)
			result = append(result, api_models.ValidationError{
				Field:   fieldName,
				Message: s.buildFieldMessage(lang, fieldName, fe),
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
		msg := s.translate(lang, "validation_type_mismatch", map[string]string{"field": field, "param": unmarshalTypeErr.Type.String()})
		return []api_models.ValidationError{{
			Field:   field,
			Message: msg,
			Tag:     "type_mismatch",
			Value:   unmarshalTypeErr.Value,
		}}
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		msg := s.translate(lang, "validation_json_syntax", map[string]string{"param": fmt.Sprintf("%d", syntaxErr.Offset)})
		return []api_models.ValidationError{{
			Field:   "body",
			Message: msg,
			Tag:     "json_syntax",
		}}
	}

	if errors.Is(err, io.EOF) {
		msg := s.translate(lang, "validation_body_empty", nil)
		return []api_models.ValidationError{{
			Field:   "body",
			Message: msg,
			Tag:     "required",
		}}
	}

	msg := s.translate(lang, "validation_parse_error", map[string]string{"field": err.Error()})
	return []api_models.ValidationError{{
		Field:   "body",
		Message: msg,
		Tag:     "parse_error",
	}}
}

// ──────────────────────────────────────────────
// Internal Helpers
// ──────────────────────────────────────────────

func (s *service) translate(lang, key string, params map[string]string) string {
	if s.Input.Services.Multilang != nil {
		return s.Input.Services.Multilang.GetMessage(lang, key, params)
	}
	// Fallback when multilang not configured
	return key
}

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

func (s *service) buildFieldMessage(lang, field string, fe validator.FieldError) string {
	param := fe.Param()
	tag := fe.Tag()
	isStr := fe.Kind() == reflect.String

	key := ""
	params := map[string]string{"field": field, "param": param}

	switch tag {
	case "required", "required_if", "required_unless", "required_with", "required_without",
		"required_with_all", "required_without_all":
		key = "validation_required"
	case "min":
		if isStr {
			key = "validation_min_str"
		} else {
			key = "validation_min"
		}
	case "max":
		if isStr {
			key = "validation_max_str"
		} else {
			key = "validation_max"
		}
	case "len":
		if isStr {
			key = "validation_len_str"
		} else {
			key = "validation_len"
		}
	case "eq":
		key = "validation_eq"
	case "ne":
		key = "validation_ne"
	case "gt":
		key = "validation_gt"
	case "gte":
		key = "validation_gte"
	case "lt":
		key = "validation_lt"
	case "lte":
		key = "validation_lte"
	case "oneof":
		key = "validation_oneof"
	case "eqfield":
		key = "validation_eqfield"
	case "nefield":
		key = "validation_nefield"
	case "gtfield":
		key = "validation_gtfield"
	case "ltfield":
		key = "validation_ltfield"
	case "email":
		key = "validation_email"
	case "url":
		key = "validation_url"
	case "uri":
		key = "validation_uri"
	case "uuid", "uuid3", "uuid4", "uuid5":
		key = "validation_uuid"
	case "ip":
		key = "validation_ip"
	case "ipv4":
		key = "validation_ipv4"
	case "ipv6":
		key = "validation_ipv6"
	case "json":
		key = "validation_json"
	case "jwt":
		key = "validation_jwt"
	case "e164":
		key = "validation_e164"
	case "alpha":
		key = "validation_alpha"
	case "alphanum":
		key = "validation_alphanum"
	case "alphanumunicode":
		key = "validation_alphanumunicode"
	case "numeric":
		key = "validation_numeric"
	case "boolean":
		key = "validation_boolean"
	case "lowercase":
		key = "validation_lowercase"
	case "uppercase":
		key = "validation_uppercase"
	case "contains":
		key = "validation_contains"
	case "excludes":
		key = "validation_excludes"
	case "startswith":
		key = "validation_startswith"
	case "endswith":
		key = "validation_endswith"
	case "datetime":
		key = "validation_datetime"
	default:
		key = "validation_default"
		params["param"] = tag
	}
	return s.translate(lang, key, params)
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
