package apperrors

import (
	"errors"
	"testing"
)

func TestErrBadReq(t *testing.T) {
	if ErrBadReq == nil {
		t.Fatal("ErrBadReq should not be nil")
	}
	if ErrBadReq.Error() != "bad request" {
		t.Errorf("ErrBadReq = %q; want bad request", ErrBadReq.Error())
	}
}

func TestErrNotFound(t *testing.T) {
	if ErrNotFound == nil {
		t.Fatal("ErrNotFound should not be nil")
	}
	if ErrNotFound.Error() != "not found" {
		t.Errorf("ErrNotFound = %q; want not found", ErrNotFound.Error())
	}
}

func TestErrInvalidInput(t *testing.T) {
	if ErrInvalidInput == nil {
		t.Fatal("ErrInvalidInput should not be nil")
	}
	if ErrInvalidInput.Error() != "invalid input" {
		t.Errorf("ErrInvalidInput = %q; want invalid input", ErrInvalidInput.Error())
	}
}

func TestErrConflict(t *testing.T) {
	if ErrConflict == nil {
		t.Fatal("ErrConflict should not be nil")
	}
	if ErrConflict.Error() != "conflict" {
		t.Errorf("ErrConflict = %q; want conflict", ErrConflict.Error())
	}
}

func TestErrInternalError(t *testing.T) {
	if ErrInternalError == nil {
		t.Fatal("ErrInternalError should not be nil")
	}
	if ErrInternalError.Error() != "internal error" {
		t.Errorf("ErrInternalError = %q; want internal error", ErrInternalError.Error())
	}
}

func TestErrorsAreDistinct(t *testing.T) {
	errs := []error{ErrBadReq, ErrNotFound, ErrInvalidInput, ErrConflict, ErrInternalError}
	for i, a := range errs {
		for j, b := range errs {
			if i != j && errors.Is(a, b) {
				t.Errorf("error %d and %d should be distinct", i, j)
			}
		}
	}
}
