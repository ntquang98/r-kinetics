package common

import "github.com/ntquang98/go-rkinetics-service/internal/pkg/encoding"

type CommonResponse struct {
	Status    string            `json:"status"`
	Message   string            `json:"message"`
	ErrorCode string            `json:"errorCode,omitempty"`
	Data      any               `json:"data,omitempty"`
	Errors    any               `json:"errors,omitempty"`
	Total     int64             `json:"total,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type Response interface {
	String() string
	Bytes() []byte
	Interface() *CommonResponse
}

// APIResponse This is response object with JSON format
type APIResponse[T any] struct {
	Status    string            `json:"status"`
	Message   string            `json:"message"`
	ErrorCode string            `json:"errorCode,omitempty"`
	Data      []T               `json:"data,omitempty"`
	Errors    any               `json:"errors,omitempty"`
	Total     int64             `json:"total,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type APIResponseSingle[T any] struct {
	Status    string            `json:"status"`
	Message   string            `json:"message"`
	ErrorCode string            `json:"errorCode,omitempty"`
	Data      T                 `json:"data,omitempty"`
	Errors    any               `json:"errors,omitempty"`
	Total     int64             `json:"total,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// statusEnum ...
type statusEnum struct {
	Ok           string
	Error        string
	Invalid      string
	NotFound     string
	Forbidden    string
	Existed      string
	Unauthorized string
}

// APIStatus Published enum
var APIStatus = &statusEnum{
	Ok:           "OK",
	Error:        "ERROR",
	Invalid:      "INVALID",
	NotFound:     "NOT_FOUND",
	Forbidden:    "FORBIDDEN",
	Existed:      "EXISTED",
	Unauthorized: "UNAUTHORIZED",
}

func (r *APIResponse[T]) Ok() bool {
	return r.Status == APIStatus.Ok
}

func (r *APIResponse[T]) Invalid() bool {
	return r.Status == APIStatus.Invalid
}

func (r *APIResponse[T]) NotFound() bool {
	return r.Status == APIStatus.NotFound
}

func (r *APIResponse[T]) Unauthorized() bool {
	return r.Status == APIStatus.Unauthorized
}

func (r *APIResponse[T]) Existed() bool {
	return r.Status == APIStatus.Existed
}

func (r *APIResponse[T]) Forbidden() bool {
	return r.Status == APIStatus.Forbidden
}

func (r *APIResponse[T]) Error() bool {
	return r.Status == APIStatus.Error
}

func (r *APIResponse[T]) String() string {
	b, _ := encoding.MarshalToString(r)
	return b
}

func (r *APIResponse[T]) Bytes() []byte {
	b, _ := encoding.Marshal(r)
	return b
}

func (r *APIResponse[T]) Interface() *CommonResponse {
	var _r = &CommonResponse{
		Status:    r.Status,
		Message:   r.Message,
		ErrorCode: r.ErrorCode,
		Errors:    r.Errors,
		Total:     r.Total,
		Headers:   r.Headers,
		Data:      r.Data,
	}
	return _r
}
