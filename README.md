# errors

A lightweight Go utility package that builds on the standard errors with ergonomic helpers for:

- Creating descriptive errors with or without stack traces
- Wrapping errors with context (formatted or plain)
- Preserving and inspecting stack traces across error chains
- Working with predefined sentinel errors (HTTP-like categories)
- Using standard errors helpers (Is/As/Unwrap) via re-exports

Note: The package name is errors, so it’s recommended to import it with an alias to avoid confusion with the standard library.

```go
import (
errs "github.com/ceearrashee/errors"
)
```

## Installation

```bash
go get github.com/ceearrashee/errors
```

## Quick start

```go
package main

import (
	"fmt"
	"os"

	errs "github.com/ceearrashee/errors"
)

func readConfig(path string) error {
	// Create a basic error with description and capture a stack
	if _, err := os.ReadFile(path); err != nil {
		return errs.Wrap(err, "reading config")
		// or: return errs.NewWithStack("config read failed")
	}
	return nil
}

func main() {
	if err := readConfig("./missing.yaml"); err != nil {
		// Find the last framework error with a stack and print it
		if e := errs.FindOriginalErrorWithStack(err); e != nil {
			fmt.Println("Error:", e.Error())
			fmt.Println("Stack trace:")
			for _, f := range e.GetCallStack() {
				fmt.Println(f)
			}
		} else {
			fmt.Println("Error:", err)
		}
	}
}
```

## Features and API highlights

- Constructors
    - `errs.New(description string) error` — simple error with description
    - `errs.NewWithStack(description string) error` — error with captured stack
    - `errs.Newf(format string, args ...any) *errs.Error` — formatted description returning the concrete type

- Wrapping helpers (nil-safe: return nil if original err is nil)
    - `errs.Wrap(err error, description string) error`
    - `errs.Wrapf(err error, format string, args ...any) error`
    - `errs.WrapWithCustomErr(originalErr, wrappingErr error) error` — wraps with a custom sentinel error
    - `errs.WrapfWithCustomErr(originalErr, wrappingErr error, format string, args ...any) error`

- Stack utilities
    - `errs.FindOriginalErrorWithStack(err error) *errs.Error` — the last framework error in the chain that has a stack
    - `errs.FindFirstErrorWithStack(err error) error` — the first framework error in the chain
    - `(*errs.Error).GetCallStack() []string` — formatted stack frames
    - `errs.AddCustomCallStack(err error, callStack *errs.Stack) error` — attach a precomputed stack (advanced)

- Standard helpers re-exported
    - `errs.Is`, `errs.As`, `errs.Unwrap` are thin wrappers around `errors.Is/As/Unwrap`
    - `errs.Errorf` is an alias of `fmt.Errorf`

## Working with predefined errors

This package ships common sentinel errors, useful for categorizing failures (HTTP-like semantics):

- `errs.ErrBadRequest` (400)
- `errs.ErrUnauthorized` (401)
- `errs.ErrRegistrationRequired` (401)
- `errs.ErrPaymentError` (402)
- `errs.ErrForbiddenAction` (403)
- `errs.ErrNotFound` (404)
- `errs.ErrConflict` (409)
- `errs.ErrPreconditionFailed` (412)
- `errs.ErrValidation` (422)
- `errs.ErrInternalServerError` (500)

Typical usage:

```go
if err == sql.ErrNoRows {
return errs.WrapWithCustomErr(err, errs.ErrNotFound)
}

// or with a message
return errs.WrapfWithCustomErr(err, errs.ErrValidation, "invalid input: %s", field)
```

And checking downstream:

```go
if errs.Is(err, errs.ErrNotFound) {
// map to HTTP 404 or similar handling
}

// If you need the first predefined error from the chain
var e *errs.Error
if errs.As(err, &e) {
fmt.Println("predefined:", e.GetOriginalPredefinedError())
}
```

## More examples

- Simple wrapping with formatted context:

```go
if err := doThing(); err != nil {
return errs.Wrapf(err, "doing thing %q", id)
}
```

- Creating a formatted error directly:

```go
return errs.Newf("user %d not found", userID)
```

- Printing the last captured stack from an error chain:

```go
if e := errs.FindOriginalErrorWithStack(err); e != nil {
for _, frame := range e.GetCallStack() {
fmt.Println(frame)
}
}
```

## Best practices

- Prefer `Wrap/Wrapf` to add context so the original cause is preserved.
- For categorization and cross-boundary handling, wrap with predefined sentinels and check via `errs.Is`.
- Always import with an alias (e.g., `errs`) to avoid confusion with the standard library `errors` package.
- Wrapping helpers are nil-safe; returning nil if the input error is nil helps reduce boilerplate.

## Compatibility

- Fully compatible with Go’s error interfaces and `errors.Is/As/Unwrap`.
- The library captures a call stack when you create or wrap using provided helpers.

## Version and requirements

- Module: `github.com/ceearrashee/errors`
- Go: 1.20+ recommended (module declares 1.24)

## License

This project is licensed under the terms of the MIT License. See the [LICENSE](./LICENSE) file for details.
