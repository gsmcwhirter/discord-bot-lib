package session

import "github.com/gsmcwhirter/go-util/v8/errors"

// ErrBadData is the error returned when upserting State elements if data is an incorrect type or invalid value
var ErrBadData = errors.New("bad data")

// ErrMissingData is the error returned when upserting State elements is missing required fields
var ErrMissingData = errors.New("missing data")

// ErrNotFound is the error returned when upserting State elements and the element to update is not found
var ErrNotFound = errors.New("not found")
