package e

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`

	err error
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", h.Code, h.Message)
}

func (h *HTTPError) Unwrap() error {
	return h.err
}

func NewError(opts ...HTTPErrorOption) *HTTPError {
	err := &HTTPError{
		Code:      http.StatusInternalServerError,
		Message:   "",
		Timestamp: time.Now(),
		err:       nil,
	}

	applyOptions(err, opts...)

	return err
}

func Internal(opts ...HTTPErrorOption) *HTTPError {
	err := NewError(
		WithStatusCode(http.StatusInternalServerError),
		WithMessage("internal server error"),
	)

	applyOptions(err, opts...)

	return err
}

func Authorization(opts ...HTTPErrorOption) *HTTPError {
	err := NewError(
		WithStatusCode(http.StatusForbidden),
		WithMessage("authorization error"),
	)

	applyOptions(err, opts...)

	return err
}

func BadRequest(opts ...HTTPErrorOption) *HTTPError {
	err := NewError(
		WithStatusCode(http.StatusBadRequest),
		WithMessage("invalid credentials"),
	)

	applyOptions(err, opts...)

	return err
}

func NotFound(opts ...HTTPErrorOption) *HTTPError {
	err := NewError(
		WithStatusCode(http.StatusNotFound),
		WithMessage("not found"),
	)

	applyOptions(err, opts...)

	return err
}

func Forbidden(opts ...HTTPErrorOption) *HTTPError {
	err := NewError(
		WithStatusCode(http.StatusForbidden),
		WithMessage("forbidden"),
	)

	applyOptions(err, opts...)

	return err
}

type HTTPErrorOption func(e *HTTPError)

func WithStatusCode(code int) HTTPErrorOption {
	return func(e *HTTPError) {
		e.Code = code
	}
}

func WithMessage(message string) HTTPErrorOption {
	return func(e *HTTPError) {
		e.Message = message
	}
}

func WithError(err error) HTTPErrorOption {
	return func(e *HTTPError) {
		e.err = err
	}
}

func applyOptions(err *HTTPError, opts ...HTTPErrorOption) {
	for _, opt := range opts {
		opt(err)
	}
}
