package err

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	convert := Convert(err)
	w.WriteHeader(err2code(convert))
	json.NewEncoder(w).Encode(convert)
}

func err2code(err *Status) int {
	return err.Code()
}

var Err404 = NewFromCode(http.StatusNotFound)
var Err401 = NewFromCode(http.StatusUnauthorized)
var Err422 = NewFromCode(http.StatusUnprocessableEntity)

type err struct {
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"errors,omitempty"`
}

func (e err) String() string {
	return e.Message
}

type Status struct {
	*err
}

func (e Status) String() string {
	return e.err.Message
}

func (e Status) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s", e.Code(), e.Message())
}

func (e Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.err)
}

func (e Status) Code() int {
	return e.err.Code
}

func (e Status) Message() string {
	return e.err.Message
}

func (e Status) Err() error {
	if e.Code() <= 400 {
		return nil
	}
	return e
}

func (e *Status) F(format string, a ...interface{}) *Status {
	e.err.Message = fmt.Sprintf(format, a...)
	return e
}

func (e *Status) Data(data interface{}) *Status {
	e.err.Data = data
	return e
}

// New returns a Status representing c and msg.
func New(c int, msg string) *Status {
	return &Status{err: &err{Code: c, Message: msg}}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(c int, format string, a ...interface{}) *Status {
	return New(c, fmt.Sprintf(format, a...))
}

// Error returns an error representing c and msg.  If c is OK, returns nil.
func Error(c int, msg string) error {
	return New(c, msg).Err()
}

// Errorf returns Error(c, fmt.Sprintf(format, a...)).
func Errorf(c int, format string, a ...interface{}) error {
	return Error(c, fmt.Sprintf(format, a...))
}

func NewFromCode(code int) *Status {
	return New(code, http.StatusText(code))
}

func Convert(err error) *Status {
	switch err.(type) {
	case *Status:
		return err.(*Status)
	case validator.ValidationErrors:
		return NewFromCode(http.StatusUnprocessableEntity).Data(err)
	default:
		return NewFromCode(500).Data(err)
	}
}
