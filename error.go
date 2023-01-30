package library

import "fmt"

//go:generate stringer -type=ErrorType
type ErrorType int

const (
	_ ErrorType = iota
	Unknown
	DatabaseError
	BadInput
	InvalidSettings
	Timeout
)

type Error struct {
	Type   ErrorType
	Desc   string
	Actual error
}

func (e *Error) Unwrap() error {
	return e.Actual
}

func (e *Error) Error() string {
	return fmt.Sprintf("(%s) %s: %s", e.Type, e.Desc, e.Actual)
}
