package errors_test

import (
	"fmt"
	"testing"

	"github.com/ceearrashee/errors"
)

// TestFindPublicErrorSearchesTheWholeTree verifies that a public error is found inside a join and
// inside a custom wrap, and that the outermost or leftmost one wins when there are several.
func TestFindPublicErrorSearchesTheWholeTree(t *testing.T) {
	public := errors.WrapPublic(errors.ErrBadRequest, "Quantity must be at least 1")
	other := errors.WrapPublic(errors.ErrConflict, "The order is already paid")
	closeErr := errors.New("closing the upload")

	for name, tc := range map[string]struct {
		err  error
		want string
	}{
		"behind a wrap in a join": {
			err:  errors.Join(errors.Wrapf(public, "creating order %d", 7), closeErr),
			want: "Quantity must be at least 1",
		},
		"in the second member of a join": {
			err:  errors.Join(closeErr, public),
			want: "Quantity must be at least 1",
		},
		"the first member of a join wins": {
			err:  errors.Join(public, other),
			want: "Quantity must be at least 1",
		},
		"under a custom wrapping error": {
			err:  errors.WrapfWithCustomErr(public, errors.ErrInternalServerError, "saving order %d", 7),
			want: "Quantity must be at least 1",
		},
		"under a multi-%w wrap": {
			err:  fmt.Errorf("%w: %w", closeErr, public),
			want: "Quantity must be at least 1",
		},
		"the outermost public error wins": {
			err:  errors.WrapPublic(errors.Wrapf(public, "creating order %d", 7), "The order cannot be created"),
			want: "The order cannot be created",
		},
	} {
		t.Run(name, func(t *testing.T) {
			found := errors.FindPublicError(tc.err)

			var frameworkErr *errors.Error
			if !errors.As(found, &frameworkErr) {
				t.Fatalf("expected a public error in %v, got %v", tc.err, found)
			}

			if frameworkErr.Description != tc.want {
				t.Errorf("Description = %q, want %q", frameworkErr.Description, tc.want)
			}
		})
	}
}

// memberBehindAs exposes its member only through an As method, as a go-multierror chain does.
type memberBehindAs struct {
	member error
}

func (m memberBehindAs) Error() string {
	return m.member.Error()
}

func (m memberBehindAs) As(target any) bool {
	return errors.As(m.member, target)
}

// TestFindPublicErrorUsesAsMethods verifies that a public error reachable only through an error's own
// As method is still found.
func TestFindPublicErrorUsesAsMethods(t *testing.T) {
	err := errors.Wrapf(memberBehindAs{member: errors.WrapPublic(errors.ErrConflict, "The order is already paid")},
		"paying order %d", 7)

	var frameworkErr *errors.Error
	if !errors.As(errors.FindPublicError(err), &frameworkErr) {
		t.Fatalf("expected a public error in %v", err)
	}

	if frameworkErr.Description != "The order is already paid" {
		t.Errorf("Description = %q, want %q", frameworkErr.Description, "The order is already paid")
	}
}

// TestFindPublicErrorIgnoresATreeWithoutOne verifies that a join or a custom wrap with no public
// error anywhere yields nil.
func TestFindPublicErrorIgnoresATreeWithoutOne(t *testing.T) {
	for name, err := range map[string]error{
		"join":        errors.Join(errors.Wrap(errors.ErrBadRequest, "parsing"), errors.New("closing")),
		"custom wrap": errors.WrapfWithCustomErr(errors.New("timeout"), errors.ErrInternalServerError, "calling"),
	} {
		t.Run(name, func(t *testing.T) {
			if found := errors.FindPublicError(err); found != nil {
				t.Errorf("expected nil, got %v", found)
			}
		})
	}
}
