package errors_test

import (
	"testing"

	"github.com/ceearrashee/errors"
)

// TestWrapfWithCustomErrPreservesBothSentinels verifies that after WrapfWithCustomErr,
// errors.Is returns true for both the wrapping sentinel and the original error.
func TestWrapfWithCustomErrPreservesBothSentinels(t *testing.T) {
	original := errors.New("original sentinel")
	wrapper := errors.New("wrapper sentinel")

	result := errors.WrapfWithCustomErr(original, wrapper, "context: %s", "value")

	if !errors.Is(result, wrapper) {
		t.Error("errors.Is(result, wrapper) must be true")
	}

	if !errors.Is(result, original) {
		t.Error("errors.Is(result, original) must be true")
	}
}

// TestWrapWithCustomErrPreservesBothSentinels verifies the same guarantee for WrapWithCustomErr.
func TestWrapWithCustomErrPreservesBothSentinels(t *testing.T) {
	original := errors.New("original sentinel")
	wrapper := errors.New("wrapper sentinel")

	result := errors.WrapWithCustomErr(original, wrapper)

	if !errors.Is(result, wrapper) {
		t.Error("errors.Is(result, wrapper) must be true")
	}

	if !errors.Is(result, original) {
		t.Error("errors.Is(result, original) must be true")
	}
}

// TestWrapfWithCustomErrNilOriginal verifies nil propagation.
func TestWrapfWithCustomErrNilOriginal(t *testing.T) {
	wrapper := errors.New("wrapper sentinel")

	if result := errors.WrapfWithCustomErr(nil, wrapper, "context"); result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

// TestWrapfWithCustomErrMessage verifies the formatted description is preserved and Error()
// returns the full chain (description + wrapped errors).
func TestWrapfWithCustomErrMessage(t *testing.T) {
	original := errors.New("original")
	wrapper := errors.ErrForbiddenAction

	result := errors.WrapfWithCustomErr(original, wrapper, "failed to make %q", "some_event")

	if result == nil {
		t.Fatal("result must not be nil")
	}

	// Error() includes the full chain: description + underlying error message.
	got := result.Error()
	want := "failed to make \"some_event\": forbidden: original"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
