package errors

import (
	"fmt"
	"runtime"
)

type (
	// Error represents an error structure that includes a description and an error string.
	// It is useful for providing additional context or user-friendly messages alongside the error text.
	Error struct {
		Description string
		error       error
		stack       *Stack
	}

	// Stack represents a slice of uintptrs, typically used to store function call stack pointers.
	Stack []uintptr
)

// Format customizes the formatted output of an Error instance.
//
// Parameters:
//   - f: the formatter state used for custom formatting
//   - _: the rune specifying the format verb (unused)
//
// Returns: none (writes the formatted description to f)
func (e *Error) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprintf(f, "%s", e.Description) //nolint:errcheck,revive
}

// Newf creates a new Error instance with a formatted description.
//
// Parameters:
//   - formatedDescription: a string with formatting placeholders for the description
//   - args: values to replace the placeholders in the formatted description
//
// Returns:
//   - *Error: a pointer to the newly created Error instance with the formatted description set.
func Newf(formatedDescription string, args ...any) *Error {
	return &Error{
		Description: fmt.Sprintf(formatedDescription, args...),
	}
}

// Error returns the error message, combining the description and underlying error if present.
//
// Returns:
//   - string: the error message, formatted as a string.
func (e *Error) Error() string {
	if e.error == nil {
		return e.Description
	}

	return fmt.Sprintf("%s: %s", e.Description, e.error.Error())
}

// Message returns the description if set; otherwise, it returns the underlying error's message.
//
// Returns:
//   - string: the error message, formatted as a string.
func (e *Error) Message() string {
	if e.Description == "" {
		return e.error.Error()
	}

	return e.Description
}

// GetOriginalErrorMessage returns the deepest error message in the error chain,
// optionally prefixed by the error's description.
//
// Returns:
//   - string: the error message, formatted as a string.
func (e *Error) GetOriginalErrorMessage() string {
	var originalErr error
	for err := Unwrap(e.error); err != nil; err = Unwrap(err) {
		originalErr = err
	}

	if e.Description == "" {
		if originalErr != nil {
			return originalErr.Error()
		}

		return e.error.Error()
	}

	errToUse := e.error

	if originalErr != nil {
		errToUse = originalErr
	}

	return fmt.Sprintf("%s: %s", e.Description, errToUse.Error())
}

// Wrap adds context to an existing error using the Error's description.
//
// Parameters:
//   - err: the error to wrap; if nil or no description is provided, returns nil.
//
// Returns:
//   - error: a new Error instance incorporating the provided error.
//
// Errors:
//   - None directly, but may wrap any provided error with additional context.
func (e *Error) Wrap(err error) error {
	if err == nil || e.Description == "" {
		return nil
	}

	return &Error{error: err, Description: e.Description, stack: callers()}
}

// Wrapf formats and wraps an existing error with the Error's description and a custom message.
//
// Parameters:
//   - format: a format string for the error message
//   - err: the error to wrap; if nil or no description is provided, returns nil
//
// Returns:
//   - error: a new error combining the format, description, and wrapped error
//
// Errors:
//   - None directly, but wraps the provided error with formatted context, if present.
func (e *Error) Wrapf(format string, err error) error {
	if err == nil || e.Description == "" {
		return nil
	}

	er := &Error{error: err, Description: e.Description, stack: callers()}

	return fmt.Errorf(format+" :%w", er) //nolint:err113
}

// Unwrap returns the wrapped error, enabling error unwrapping in chains and supporting the errors.Unwrap interface.
func (e *Error) Unwrap() error {
	return e.error
}

// GetCallStack retrieves the function call stack associated with the error.
//
// Returns:
//   - []string: a slice of formatted call stack frames as strings, in order from most to least recent.
func (e *Error) GetCallStack() []string {
	if e == nil {
		return nil
	}

	if e.stack == nil {
		return nil
	}

	callStackFrames := make([]string, 0, 32)
	frames := runtime.CallersFrames(*e.stack)

	for {
		frame, more := frames.Next()
		if frame.Function == "unknown" {
			break
		}

		callStackFrames = append(callStackFrames, fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line))

		if !more {
			break
		}
	}

	return callStackFrames
}

func callers() *Stack {
	const depth = 32

	var pcs [depth]uintptr

	n := runtime.Callers(3, pcs[:]) //nolint:mnd

	var st Stack = pcs[0:n]

	return &st
}

// GetOriginalPredefinedError retrieves the first predefined error in the error chain if any exist.
//
// Returns:
//   - error: the first predefined error in the chain, or the original error if no predefined error is found.
func (e *Error) GetOriginalPredefinedError() error {
	var predefinedErr = e.error

	for err := Unwrap(e.error); err != nil; err = Unwrap(err) {
		switch {
		case Is(err, ErrBadRequest),
			Is(err, ErrUnauthorized),
			Is(err, ErrRegistrationRequired),
			Is(err, ErrPaymentError),
			Is(err, ErrForbiddenAction),
			Is(err, ErrNotFound),
			Is(err, ErrConflict),
			Is(err, ErrPreconditionFailed),
			Is(err, ErrValidation),
			Is(err, ErrInternalServerError):
			predefinedErr = err
		default:
			return predefinedErr
		}
	}

	return predefinedErr
}
