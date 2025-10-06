package errors

import (
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
)

// Error is wrapped error with code and inner error
type Error interface {
	error
	GetCode() int
	Error() string
}

// AppError ...
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Inner   error  `json:"-"`
}

// Error ...
func (e *AppError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%d: %s: %v", e.Code, e.Message, e.Inner)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// GetCode ...
func (e *AppError) GetCode() int {
	return e.Code
}

// NewInvalidError ...
func NewInvalidError(params ...string) *AppError {
	return &AppError{
		Code:    InvalidCode,
		Message: fmt.Sprintf("invalid param: %v", params),
	}
}

// NewNotFoundError ...
func NewNotFoundError(resourceType, content string) *AppError {
	return &AppError{
		Code:    NotFoundCode,
		Message: fmt.Sprintf("%s %s not found", resourceType, content),
	}
}

// NewCannotExecError ...
func NewCannotExecError(msg string) *AppError {
	return &AppError{
		Code:    CannotExecCode,
		Message: msg,
	}
}

// NewInternalError ...
func NewInternalError(err error) *AppError {
	return &AppError{
		Code:    InternalCode,
		Message: "internal system error",
		Inner:   err,
	}
}

// NewHertzRouteNotFoundError ...
func NewHertzRouteNotFoundError(ctx *app.RequestContext) *AppError {
	return &AppError{
		Code:    RouteNotFoundCode,
		Message: fmt.Sprintf("route %s not found", ctx.Request.Path()),
	}
}

// NewHertzBindError ...
func NewHertzBindError(err error) *AppError {
	return &AppError{
		Code:    BindErrorCode,
		Message: "hertz bind error",
		Inner:   err,
	}
}

// IsCode ...
func IsCode(err error, code int) bool {
	var appError *AppError
	return errors.As(err, &appError) && appError.GetCode() == code
}
