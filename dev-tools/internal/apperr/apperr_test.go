package apperr

import (
	"errors"
	"fmt"
	"testing"
)

func TestError_WithFieldAndWrapReturnCopies(t *testing.T) {
	orig := New(Invalid, "project.name_required")
	cause := errors.New("causa")

	withField := orig.WithField("name")
	wrapped := orig.Wrap(cause)

	if orig.Field != "" || orig.Err != nil {
		t.Errorf("o original foi alterado: %+v", orig)
	}
	if withField == orig || withField.Field != "name" || withField.Kind != Invalid || withField.Code != "project.name_required" {
		t.Errorf("WithField = %+v", withField)
	}
	if wrapped == orig || wrapped.Err != cause || wrapped.Code != "project.name_required" {
		t.Errorf("Wrap = %+v", wrapped)
	}
}

func TestError_Is(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"mesmo código com campo e causa", ErrInvalidBody.WithField("x").Wrap(errors.New("c")), ErrInvalidBody, true},
		{"mesmo código embrulhado por fmt.Errorf", fmt.Errorf("ctx: %w", ErrTokenInvalid), ErrTokenInvalid, true},
		{"código diferente", ErrInvalidBody, ErrTokenInvalid, false},
		{"alvo que não é *Error", ErrInvalidBody, errors.New("request.invalid_body"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("errors.Is = %v, quer %v", got, tt.want)
			}
		})
	}
}

func TestError_ErrorAndUnwrap(t *testing.T) {
	cause := errors.New("disk full")
	err := New(Internal, "internal.unexpected").Wrap(cause)

	if got, want := err.Error(), "internal.unexpected: disk full"; got != want {
		t.Errorf("Error() = %q, quer %q", got, want)
	}
	if !errors.Is(err, cause) {
		t.Error("errors.Is não encontrou a causa via Unwrap")
	}
	if got, want := New(Internal, "internal.unexpected").Error(), "internal.unexpected"; got != want {
		t.Errorf("Error() sem causa = %q, quer %q", got, want)
	}
}

func TestCommonErrors(t *testing.T) {
	tests := []struct {
		err  *Error
		kind Kind
		code string
	}{
		{ErrInvalidBody, BadRequest, "request.invalid_body"},
		{ErrTokenInvalid, Unauthorized, "auth.token_invalid"},
		{ErrDatabaseUnavailable, Unavailable, "database.unavailable"},
		{ErrInternal, Internal, "internal.unexpected"},
	}
	for _, tt := range tests {
		if tt.err.Kind != tt.kind || tt.err.Code != tt.code {
			t.Errorf("%s: Kind=%s Code=%s, quer %s/%s", tt.code, tt.err.Kind, tt.err.Code, tt.kind, tt.code)
		}
	}
}
