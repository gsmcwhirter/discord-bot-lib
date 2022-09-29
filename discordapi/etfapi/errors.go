package etfapi

import "github.com/gsmcwhirter/go-util/v10/errors"

// ErrBadTarget is the error returned when trying to convert
// an Element to an incorrect type
var ErrBadTarget = errors.New("bad element unmarshal target")

// ErrBadPayload is the error returned when attempting to unmarshal
// an etf payload fails due to bad formatting
var ErrBadPayload = errors.New("bad payload format")

// ErrBadFieldType is the error returned when attempting to unmarshal
// an etf payload and a field is an incorrect type (like non-string-like map keys)
var ErrBadFieldType = errors.New("bad field type")

// ErrBadMarshalData is the error returned when attempting to marshal
// an etf payload to []byte and the data in an Element doesn't match the Code
var ErrBadMarshalData = errors.New("bad marshal data")

// ErrBadElementData is the error returned when attempting to create a new element
// but the data provided does not match the Code
var ErrBadElementData = errors.New("bad element data")

// ErrOutOfBounds is the error returned when integer values are out of the bounds
// of the type code
var ErrOutOfBounds = errors.New("int value out of bounds")

// ErrBadParity is the error returned when a list that should be even parity is not
// (usually when trying to deal with Map Elements)
var ErrBadParity = errors.New("non-even list parity")
