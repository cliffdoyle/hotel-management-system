package validator

import (
	"errors"
	"regexp"
)

var (
	EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// Validator holds our validation errors.
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the Errors map is empty.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

var ErrInvalidCredentials = errors.New("invalid credentials")
