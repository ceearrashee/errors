package errors_test

import (
	"fmt"
	"testing"

	"github.com/ceearrashee/errors"
)

// ExampleWrapPublic demonstrates wrapping an error with a public message.
func ExampleWrapPublic() {
	internalErr := errors.New("database connection failed")
	publicErr := errors.WrapPublic(internalErr, "Unable to connect to service")

	var frameworkErr *errors.Error
	if errors.As(publicErr, &frameworkErr) {
		fmt.Println("Is public:", frameworkErr.IsPublic())
		fmt.Println("Message:", frameworkErr.Message())
	}

	// Output:
	// Is public: true
	// Message: Unable to connect to service
}

// ExampleWrapfPublic demonstrates wrapping an error with a formatted public message.
func ExampleWrapfPublic() {
	internalErr := errors.New("SQL error: syntax error near line 42")
	publicErr := errors.WrapfPublic(internalErr, "Failed to process request #%d", 12345)

	var frameworkErr *errors.Error
	if errors.As(publicErr, &frameworkErr) {
		fmt.Println("Is public:", frameworkErr.IsPublic())
		fmt.Println("Message:", frameworkErr.Message())
	}

	// Output:
	// Is public: true
	// Message: Failed to process request #12345
}

// ExampleFindPublicError demonstrates finding the first public error in an error chain.
func ExampleFindPublicError() {
	// Create a chain of errors with one public error
	internalErr := errors.New("internal database error")
	wrappedErr := errors.Wrap(internalErr, "query failed")
	publicErr := errors.WrapPublic(wrappedErr, "Service temporarily unavailable")
	finalErr := errors.Wrap(publicErr, "handler error")

	// Find the public error
	var frameworkErr *errors.Error
	if errors.As(errors.FindPublicError(finalErr), &frameworkErr) {
		fmt.Println("Found public error:", frameworkErr.Message())
	}

	// Output:
	// Found public error: Service temporarily unavailable
}

// TestFindPublicErrorNotFound tests that FindPublicError returns nil when the chain has no public error.
func TestFindPublicErrorNotFound(t *testing.T) {
	for name, err := range map[string]error{
		"nil input":       nil,
		"no public error": errors.Wrap(errors.New("internal"), "wrapped"),
	} {
		t.Run(name, func(t *testing.T) {
			if found := errors.FindPublicError(err); found != nil {
				t.Errorf("expected nil, got %v", found)
			}
		})
	}
}
