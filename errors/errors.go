package errors

import (
	"fmt"
)

// FullReinitializationError depicts that an error occured due to
// the fact that the browser context isn't even in the correct location
// such that re-clicking the 'Search for Class Sections' link can be clicked
type FullReinitializationError struct {
	msg string
}

// FullReinitializationError.Error() returns the errors message for the FullReinitializationError
func (r *FullReinitializationError) Error() string {
	return r.msg
}

// NewFullReinitializationError returns a new FullReinitializationError
func NewFullReinitializationError(msg string) error {
	return &FullReinitializationError{msg}
}

// WrapFullReinitializationError returns an error where the new
// message is combined with the old message to create the message
// for the FullReinitializationError which is returned
func WrapFullReinitializationError(err error, msg string) error {
	return &FullReinitializationError{fmt.Sprintf("%s: %s", msg, err.Error())}
}

// ReinitializationError depicts that an error occured due to
// the fact that the search area is not present
type ReinitializationError struct {
	msg string
}

// ReinitializationError.Error() returns the errors message for the ReinitializationError
func (r *ReinitializationError) Error() string {
	return r.msg
}

// NewReinitializationError returns a new ReinitializationError
func NewReinitializationError(msg string) error {
	return &ReinitializationError{msg}
}

// WrapReinitializationError returns an error where the new
// message is combined with the old message to create the message
// for the ReinitializationError which is returned
func WrapReinitializationError(err error, msg string) error {
	return &ReinitializationError{fmt.Sprintf("%s: %s", msg, err.Error())}
}
