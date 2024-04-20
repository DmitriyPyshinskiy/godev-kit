package xerrors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
)

var (
	ErrInvalidInput   = &AppError{AppCode: "INVALID_INPUT", HTTPCode: http.StatusBadRequest}
	ErrEntryNotFound  = &AppError{AppCode: "ENTRY_NOT_FOUND", HTTPCode: http.StatusNotFound}
	ErrInternalServer = &TypedAppError{ErrType: "SERVER_ERROR", AppError: &AppError{AppCode: "INTERNAL_SERVER", HTTPCode: http.StatusInternalServerError}}
)

type AppError struct {
	AppCode  string
	HTTPCode int
}

func (e *AppError) Code() string {
	return e.AppCode
}

func (e *AppError) Error() string {
	return e.AppCode
}

type TypedAppError struct {
	*AppError
	ErrType string
}

type TypedAppErrorJSON struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

func Test_XError(t *testing.T) {
	tests := []struct {
		name       string
		app        appError
		message    string
		errs       []error
		marshaller MarshalFunc
		wantDTO    any
	}{
		{
			name:    "simple error",
			app:     ErrInvalidInput,
			message: "invalid user input",
			errs:    nil,
			wantDTO: DefaultJSON{
				Code: ErrInvalidInput.AppCode,
			},
		},
		{
			name:    "with one nested error",
			app:     ErrEntryNotFound,
			message: "entry not found",
			errs:    []error{sql.ErrNoRows},
			wantDTO: DefaultJSON{
				Code: ErrEntryNotFound.AppCode,
			},
		},
		{
			name:    "with many nested errors",
			app:     ErrInternalServer,
			message: "internal error",
			errs:    []error{sql.ErrConnDone, os.ErrClosed},
			marshaller: func(e xErrorDTO) ([]byte, error) {
				var x *TypedAppError
				if errors.As(e.AppError, &x) {
					return json.Marshal(TypedAppErrorJSON{
						Code: x.Code(),
						Type: x.ErrType,
					})
				}
				return nil, ErrMarshal
			},
			wantDTO: TypedAppErrorJSON{
				Code: ErrInternalServer.AppCode,
				Type: ErrInternalServer.ErrType,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.marshaller != nil {
				SetupMarshaller(tt.marshaller)
			}

			xErr := New(tt.app, tt.message, tt.errs...)

			if !errors.Is(xErr, tt.app) {
				t.Errorf("actualError = %v, wantError %v", xErr, tt.app)
			}

			if tt.errs != nil {
				for _, err := range tt.errs {
					if !errors.Is(xErr, err) {
						t.Errorf("actualError = %v, wantError %v", xErr, err)
					}
				}
			}

			actualDTO, err := xErr.MarshalJSON()
			if err != nil {
				t.Errorf("xErr.MarshalJSON() error = %v", err)
			}

			wantDTO := marshalXErrorDTO(t, tt.wantDTO)

			if string(actualDTO) != string(wantDTO) {
				t.Errorf("actualDTO = %v, wantDTO %v", string(actualDTO), string(wantDTO))
			}
		})
	}
}

func marshalXErrorDTO(t *testing.T, xErr any) []byte {
	bytes, err := json.Marshal(xErr)
	if err != nil {
		t.Fatalf("json.Marshal(xErr) error = %v", err)
	}
	return bytes
}
