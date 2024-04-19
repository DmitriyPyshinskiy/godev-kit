package xerrors

import (
	"encoding/json"
	"errors"
	"fmt"
)

var marshaller MarshalFunc

var ErrMarshal = errors.New("error marshalling")

var DefaultMarshaller = func(e xErrorDTO) ([]byte, error) {
	var x appError
	if errors.As(e.AppError, &x) {
		return json.Marshal(DefaultJSON{Code: x.Code()})
	}
	return nil, ErrMarshal
}

func init() {
	marshaller = DefaultMarshaller
}

type MarshalFunc func(error xErrorDTO) ([]byte, error)

type DefaultJSON struct {
	Code any `json:"code"`
}

func SetupMarshaller(m MarshalFunc) {
	marshaller = m
}

type appError interface {
	Code() string
	Error() string
}

type XError[E appError] struct {
	// Application error
	app E

	// Any nested errors, may be nil
	errs []error

	// Custom message to better describing the error
	message string
}

type xErrorDTO struct {
	AppError error
	Errors   []error
	Message  string
}

func New[E appError](app E, message string, errs ...error) *XError[E] {
	e := &XError[E]{
		app:     app,
		errs:    errs,
		message: message,
	}
	return e
}

func (e XError[E]) GetApp() E {
	return e.app
}

func (e XError[E]) Error() string {
	return fmt.Sprintf("[%s]: %s", e.app.Code(), errorMessage(e.message, e.errs...))
}

func (e XError[E]) Is(target error) bool {
	// check app error
	var x appError
	if errors.As(target, &x) {
		return x.Code() == e.app.Code()
	}

	// check nested errors
	for _, err := range e.errs {
		if errors.Is(target, err) {
			return true
		}
	}

	return false
}

func (e XError[E]) MarshalJSON() ([]byte, error) {
	return marshaller(
		xErrorDTO{
			AppError: e.app,
			Errors:   e.errs,
			Message:  e.Error(),
		},
	)
}

func errorMessage(message string, errs ...error) string {
	if errs == nil {
		return message
	}

	errMsg := errors.Join(errs...).Error()

	if message == "" {
		return errMsg
	}

	return fmt.Sprintf("%s: %s", message, errMsg)
}
