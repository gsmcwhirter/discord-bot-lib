package etfapi

import "github.com/pkg/errors"

// ErrNotImplemented TODOC
var ErrNotImplemented = errors.New("not yet implemented")

// ErrBadTarget TODOC
var ErrBadTarget = errors.New("bad element unmarshal target")

// ErrBadPayload TODOC
var ErrBadPayload = errors.New("bad payload format")

// ErrBadFieldType TODOC
var ErrBadFieldType = errors.New("bad field type")

// ErrBadMarshalData TODOC
var ErrBadMarshalData = errors.New("bad marshal data")

// ErrBadElementData TODOC
var ErrBadElementData = errors.New("bad element data")

// ErrOutOfBounds TODOC
var ErrOutOfBounds = errors.New("int value out of bounds")

// ErrBadParity TODOC
var ErrBadParity = errors.New("non-even list parity")

// ErrMissingData TODOC
var ErrMissingData = errors.New("missing data")

// ErrBadData TODOC
var ErrBadData = errors.New("bad data")

// ErrNotFound TODOC
var ErrNotFound = errors.New("not found")
