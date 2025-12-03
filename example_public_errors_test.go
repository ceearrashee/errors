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
	foundErr := errors.FindPublicError(finalErr)
	if foundErr != nil {
		var frameworkErr *errors.Error
		if errors.As(foundErr, &frameworkErr) {
			fmt.Println("Found public error:", frameworkErr.Message())
		}
	}

	// Output:
	// Found public error: Service temporarily unavailable
}

// TestWrapPublic tests the WrapPublic function.
func TestWrapPublic(t *testing.T) {
	err := errors.New("internal error")
	publicErr := errors.WrapPublic(err, "public message")

	if publicErr == nil {
		t.Fatal("expected non-nil error")
	}

	var frameworkErr *errors.Error
	if !errors.As(publicErr, &frameworkErr) {
		t.Fatal("expected error to be of type *errors.Error")
	}

	if !frameworkErr.IsPublic() {
		t.Error("expected error to be marked as public")
	}

	if frameworkErr.Message() != "public message" {
		t.Errorf("expected message 'public message', got '%s'", frameworkErr.Message())
	}
}

// TestWrapfPublic tests the WrapfPublic function.
func TestWrapfPublic(t *testing.T) {
	err := errors.New("internal error")
	publicErr := errors.WrapfPublic(err, "public message: %s", "test")

	if publicErr == nil {
		t.Fatal("expected non-nil error")
	}

	var frameworkErr *errors.Error
	if !errors.As(publicErr, &frameworkErr) {
		t.Fatal("expected error to be of type *errors.Error")
	}

	if !frameworkErr.IsPublic() {
		t.Error("expected error to be marked as public")
	}

	expectedMsg := "public message: test"
	if frameworkErr.Message() != expectedMsg {
		t.Errorf("expected message '%s', got '%s'", expectedMsg, frameworkErr.Message())
	}
}

// TestFindPublicError tests the FindPublicError function.
func TestFindPublicError(t *testing.T) {
	t.Run("finds public error in chain", func(t *testing.T) {
		internalErr := errors.New("internal")
		wrappedErr := errors.Wrap(internalErr, "wrapped")
		publicErr := errors.WrapPublic(wrappedErr, "public")
		finalErr := errors.Wrap(publicErr, "final")

		found := errors.FindPublicError(finalErr)
		if found == nil {
			t.Fatal("expected to find public error")
		}

		var frameworkErr *errors.Error
		if !errors.As(found, &frameworkErr) {
			t.Fatal("expected error to be of type *errors.Error")
		}

		if frameworkErr.Message() != "public" {
			t.Errorf("expected message 'public', got '%s'", frameworkErr.Message())
		}
	})

	t.Run("returns nil when no public error exists", func(t *testing.T) {
		err := errors.New("internal")
		wrapped := errors.Wrap(err, "wrapped")

		found := errors.FindPublicError(wrapped)
		if found != nil {
			t.Error("expected nil when no public error exists")
		}
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		found := errors.FindPublicError(nil)
		if found != nil {
			t.Error("expected nil for nil input")
		}
	})
}

// TestWrapPublicWithNilError tests that WrapPublic returns nil for nil input.
func TestWrapPublicWithNilError(t *testing.T) {
	err := errors.WrapPublic(nil, "message")
	if err != nil {
		t.Error("expected nil error for nil input")
	}
}

// TestWrapfPublicWithNilError tests that WrapfPublic returns nil for nil input.
func TestWrapfPublicWithNilError(t *testing.T) {
	err := errors.WrapfPublic(nil, "message: %s", "test")
	if err != nil {
		t.Error("expected nil error for nil input")
	}
}
