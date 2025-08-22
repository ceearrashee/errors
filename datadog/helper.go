package datadog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/ceearrashee/errors"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

type (
	// RequestInfo carries optional HTTP request information for error enrichment.
	// Only Method and URI are required for basic usage.
	// Headers and Body are optional and should omit sensitive data if provided.
	RequestInfo struct {
		// Method specifies the HTTP method (e.g., GET, POST, etc.) used in the request.
		Method string `json:"method,omitempty"`
		// URI specifies the target resource's identifier in the HTTP request.
		URI string `json:"uri,omitempty"`
		// Headers contain HTTP headers associated with the request,
		// where keys are header names and values are header values.
		Headers map[string]string `json:"headers,omitempty"`
		// Body contains the HTTP request body, which may include textual or JSON data.
		Body string `json:"body,omitempty"`
	}
	// Context key type to avoid collisions.
	ctxKey int
)

const (
	requestInfoKey ctxKey = iota
)

// WithRequest attaches the provided RequestInfo to the context for further retrieval.
//
// Parameters:
//   - ctx: the parent context to derive from
//   - info: the RequestInfo to attach to the context
//
// Returns:
//   - context.Context: derived context containing the RequestInfo
func WithRequest(ctx context.Context, info RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfoKey, info)
}

// HandleError reports an error to a tracing span, adding detailed context and stack trace.
//
// Parameters:
//   - ctx: the context containing the tracing information
//   - err: the error to handle and report
//
// Behavior:
//   - Adds error details, including stack trace, to a tracing span if it's available in the given context.
//   - Tags the span with HTTP-related metadata, if present in the context.
func HandleError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	span, _ := tracer.SpanFromContext(ctx)
	if span == nil {
		return nil
	}

	defer span.Finish()

	var (
		typedError *errors.Error
		stack      string
	)

	if errors.As(err, &typedError) {
		stack = strings.Join(typedError.GetCallStack(), "\n")
	} else {
		var er error

		stack, er = buildStack(2)
		if er != nil {
			return errors.Wrapf(er, "failed to build stack trace")
		}
	}

	// Build application stack skipping helper frames.

	// Mark span as error with details compatible with DataDog UI.
	span.SetTag(ext.Error, true)
	span.SetTag(ext.ErrorMsg, err.Error())
	span.SetTag(ext.ErrorType, fmt.Sprintf("%T", err))
	span.SetTag(ext.ErrorStack, stack)
	setSpanRequestInfo(ctx, span)

	return nil
}

func setSpanRequestInfo(ctx context.Context, span *tracer.Span) {
	// Attach HTTP info if present in ctx.
	v := ctx.Value(requestInfoKey)
	if v == nil {
		return
	}

	ri, ok := v.(RequestInfo)
	if !ok {
		return
	}

	if ri.Method != "" {
		span.SetTag(ext.HTTPMethod, ri.Method)
	}

	if ri.URI != "" {
		span.SetTag(ext.HTTPURL, ri.URI)
	}

	// Compact details blob (custom tag) for extra context.
	if details := compactDetails(ri); details != "" {
		span.SetTag("error.details", details)
	}
}

func compactDetails(ri RequestInfo) string {
	extraData := make(map[string]any)
	if ri.Method != "" {
		extraData["method"] = ri.Method
	}

	if ri.URI != "" {
		extraData["uri"] = ri.URI
	}

	if len(ri.Headers) > 0 {
		extraData["headers"] = ri.Headers
	}

	if ri.Body != "" {
		// Beware of PII: caller should already have scrubbed sensitive data.
		extraData["body"] = ri.Body
	}

	if len(extraData) == 0 {
		return ""
	}

	b, err := json.Marshal(extraData)
	if err != nil {
		return ""
	}

	return string(b)
}

// buildStack renders a human-friendly call stack, skipping the first `skip` frames.
func buildStack(skip int) (string, error) {
	pcs := make([]uintptr, 64)
	n := runtime.Callers(skip, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)

	var buffer bytes.Buffer

	for {
		f, more := frames.Next()
		if f.Function != "" && !strings.HasPrefix(f.Function, "runtime.") {
			_, err := fmt.Fprintf(&buffer, "%s\n\t%s:%d\n", f.Function, f.File, f.Line)
			if err != nil {
				return "", errors.Wrapf(err, "failed to format stack trace")
			}
		}

		if !more {
			break
		}
	}

	return strings.TrimSpace(buffer.String()), nil
}
