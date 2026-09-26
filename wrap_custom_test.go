package errors_test

import (
	"testing"

	"github.com/ceearrashee/errors"
)

// TestWrapWithCustomErrPreservesBothSentinels verifies that errors.Is returns true for both
// the wrapping sentinel and the original error.
func TestWrapWithCustomErrPreservesBothSentinels(t *testing.T) {
	for name, wrap := range map[string]func(original, wrapper error) error{
		"WrapWithCustomErr": errors.WrapWithCustomErr,
		"WrapfWithCustomErr": func(original, wrapper error) error {
			return errors.WrapfWithCustomErr(original, wrapper, "context: %s", "value")
		},
	} {
		t.Run(name, func(t *testing.T) {
			original := errors.New("original sentinel")
			wrapper := errors.New("wrapper sentinel")

			result := wrap(original, wrapper)

			if !errors.Is(result, wrapper) {
				t.Error("errors.Is(result, wrapper) must be true")
			}

			if !errors.Is(result, original) {
				t.Error("errors.Is(result, original) must be true")
			}
		})
	}
}

// TestWrapNilReturnsNil verifies nil propagation for every Wrap* constructor.
func TestWrapNilReturnsNil(t *testing.T) {
	wrapper := errors.New("wrapper sentinel")

	for name, err := range map[string]error{
		"Wrap":               errors.Wrap(nil, "message"),
		"Wrapf":              errors.Wrapf(nil, "message: %s", "test"),
		"WrapPublic":         errors.WrapPublic(nil, "message"),
		"WrapfPublic":        errors.WrapfPublic(nil, "message: %s", "test"),
		"WrapWithCustomErr":  errors.WrapWithCustomErr(nil, wrapper),
		"WrapfWithCustomErr": errors.WrapfWithCustomErr(nil, wrapper, "context"),
	} {
		if err != nil {
			t.Errorf("%s: expected nil, got %v", name, err)
		}
	}
}

// TestWrapfWithCustomErrMessage verifies the formatted description is preserved and Error()
// returns the full chain (description + wrapped errors).
func TestWrapfWithCustomErrMessage(t *testing.T) {
	result := errors.WrapfWithCustomErr(errors.New("original"), errors.ErrForbiddenAction, "failed to make %q", "some_event")

	if result == nil {
		t.Fatal("result must not be nil")
	}

	got := result.Error()
	want := "failed to make \"some_event\": forbidden: original"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
