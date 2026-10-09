package domain

import "errors"

// FieldError describes a validation failure on a single input field.
type FieldError struct {
	Field  string `json:"field"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ValidationError carries one FieldError per violated input.
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

// ErrNotFound is returned when an entity does not exist.
var ErrNotFound = errors.New("not found")
