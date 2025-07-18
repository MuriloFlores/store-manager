package web_errors

import (
	"encoding/json"
	"net/http"
)

type RestErr struct {
	Message string   `json:"message"`
	Err     string   `json:"error"`
	Code    int      `json:"code"`
	Causes  []Causes `json:"causes"`
}

type Causes struct {
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func (r *RestErr) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	json.NewEncoder(w).Encode(r)
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func NewUnauthorizedRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	}
}

func NewBadRequestValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	}
}

func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	}
}

func NewConflictError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "conflict",
		Code:    http.StatusConflict,
	}
}

func NewConflictErrorWithCause(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "conflict",
		Code:    http.StatusConflict,
		Causes:  causes,
	}
}

func NewEmailNotVerified(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "EmailNotVerified",
		Code:    http.StatusForbidden,
		Causes:  causes,
	}
}

func NewRateLimitExceededError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "rate_limit_exceeded",
		Code:    http.StatusTooManyRequests,
	}
}
