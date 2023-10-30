package api

import (
	"context"
	"errors"

	"github.com/monzo/terrors"

	"github.com/goccy/go-json"

	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RenderJSON Render a helper function to render a JSON response
func RenderJSON(ctx context.Context, w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	// Headers
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(ctx))
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(js)
}

// RenderError Renders an error with some sane defaults.
func RenderError(ctx context.Context, w http.ResponseWriter, err error) {
	httpStatusCode := http.StatusInternalServerError
	code := "internal_error"
	message := "something went wrong, please try again later"

	// Check if the error is a terrors.Error
	var terror *terrors.Error
	if !errors.As(err, &terror) {
		terror = terrors.InternalService("", "something went wrong, please try again later", nil)
	}

	switch terror.Code {
	case terrors.ErrUnauthorized:
		httpStatusCode = http.StatusUnauthorized
		code = "unauthorized"
		message = terror.Message
	case terrors.ErrForbidden:
		httpStatusCode = http.StatusForbidden
		code = "forbidden"
		message = terror.Message
	case terrors.ErrNotFound:
		httpStatusCode = http.StatusNotFound
		code = "not_found"
		message = terror.Message
	case terrors.ErrPreconditionFailed:
		httpStatusCode = http.StatusBadRequest
		code = "precondition_failed"
		message = terror.Message
	case terrors.ErrBadRequest:
		httpStatusCode = http.StatusBadRequest
		code = "bad_request"
		message = terror.Message
	}

	payload := map[string]string{
		"code":    code,
		"message": message,
	}

	RenderJSON(ctx, w, httpStatusCode, payload)
}
