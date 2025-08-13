package errors

import (
	stdErrors "errors"
	"fmt"
)

var (
	// Is a wrapper for errors.Is.
	Is = stdErrors.Is
	// Unwrap is a wrapper for errors.Unwrap.
	Unwrap = stdErrors.Unwrap
	// As is a wrapper for errors.As.
	As = stdErrors.As
	// Errorf is a wrapper for fmt.Errorf.
	Errorf = fmt.Errorf //nolint:gochecknoglobals
)

// FindOriginalErrorWithStack traverses an error chain to locate the latest framework error containing a call stack.
//
// Parameters:
//   - err: the root error to search through
//
// Returns:
//   - *Error: the last framework error containing a call stack, or nil if none are found
func FindOriginalErrorWithStack(err error) *Error {
	var lastFrameworkErrWithStack *Error

	current := err

	// Traverse the entire error chain.
	for current != nil {
		var frameworkErr *Error
		if As(current, &frameworkErr) && frameworkErr.GetCallStack() != nil {
			// Found a framework error with stack, save it.
			lastFrameworkErrWithStack = frameworkErr
		}

		// Continue unwrapping.
		current = Unwrap(current)
	}

	return lastFrameworkErrWithStack
}

// FindFirstErrorWithStack traverses an error chain to locate the first framework-specific error.
//
// Parameters:
//   - err: the root error to traverse
//
// Returns:
//   - *Error: the first framework-specific error in the chain, or nil if not found
func FindFirstErrorWithStack(err error) error {
	current := err

	// Traverse the entire error chain.
	for current != nil {
		var frameworkErr *Error
		if As(current, &frameworkErr) {
			return frameworkErr
		}

		// Continue unwrapping.
		current = Unwrap(current)
	}

	return current
}

// New creates a new Error instance with the specified description.
//
// Parameters:
//   - description: a text message describing the error.
//
// Returns:
//   - error: an Error instance encapsulating the provided description.
func New(description string) error {
	return &Error{
		Description: description,
	}
}

// NewWithStack creates a new error with a description and captures the current call stack.
//
// Parameters:
//   - description: the error description
//
// Returns:
//   - error: a newly created error with stack trace included
func NewWithStack(description string) error {
	return &Error{
		Description: description,
		stack:       callers(),
	}
}

// Wrap wraps an existing error with additional context and a stack trace.
//
// Parameters:
//   - err: the original error to wrap
//   - description: a description providing context for the error
//
// Returns:
//   - error: a wrapped error with the original error, description, and stack trace, or nil if the input error is nil
func Wrap(err error, description string) error {
	if err == nil {
		return nil
	}

	return &Error{
		Description: description,
		stack:       callers(),
		error:       err,
	}
}

// Wrapf logs the given error with a formatted message and wraps the error with the same message.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return &Error{
		Description: fmt.Sprintf(format, args...),
		stack:       callers(),
		error:       err,
	}
}

// WrapfWithCustomErr creates a new Error instance by wrapping an original error with a custom error and formatted message.
//
// Parameters:
//   - originalErr: the error to wrap; if nil, the function returns nil
//   - wrappingErr: the error to use for wrapping the original error
//   - format: a format string for the custom error description
//   - args: optional arguments for formatting the description
//
// Returns:
//   - error: an Error with formatted description and wrapped errors, or nil if the original error is nil
func WrapfWithCustomErr(originalErr, wrappingErr error, format string, args ...any) error {
	if originalErr == nil {
		return nil
	}

	return &Error{
		Description: fmt.Sprintf(format, args...),
		stack:       callers(),
		error:       fmt.Errorf("%w: %v", wrappingErr, originalErr),
	}
}

// WrapWithCustomErr wraps an original error with a custom error, maintaining context and a call stack.
//
// Parameters:
//   - originalErr: the error to be wrapped
//   - wrappingErr: the error providing additional context
//
// Returns:
//   - error: a new error combining the original and custom errors, or nil if the original error is nil
func WrapWithCustomErr(originalErr, wrappingErr error) error {
	if originalErr == nil {
		return nil
	}

	return &Error{
		stack: callers(),
		error: fmt.Errorf("%w: %v", wrappingErr, originalErr),
	}
}

// AddCustomCallStack wraps the given error with a custom call stack and returns a new error that includes both.
// It preserves the original error message while providing additional call stack context useful for debugging.
//
// Parameters:
//   - err: The original error to wrap. If nil, returns nil.
//   - callStack: The custom call stack to attach to the error. Must be a valid Stack pointer.
//
// Returns:
//   - A new Error that embeds both the original error and the provided call stack.
//   - nil if the input error is nil.
func AddCustomCallStack(err error, callStack *Stack) error {
	if err == nil {
		return nil
	}

	return &Error{
		Description: err.Error(),
		stack:       callStack,
		error:       err,
	}
}
