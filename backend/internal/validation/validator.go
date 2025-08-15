package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator holds validation rules and errors
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true if there are no validation errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a given key
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message for a given key if the condition is not met
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Validation helper functions

// NotBlank checks if a value is not blank
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars checks if a string has at most n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinChars checks if a string has at least n characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// In checks if a value is in a list of permitted values
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// Matches checks if a string matches a regular expression pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all values in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}

// Common validation patterns
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PhoneRX = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

// Common validation functions

// ValidateEmail validates email format
func ValidateEmail(v *Validator, email, fieldName string) {
	v.Check(NotBlank(email), fieldName, "must be provided")
	v.Check(Matches(email, EmailRX), fieldName, "must be a valid email address")
}

// ValidatePassword validates password strength
func ValidatePassword(v *Validator, password, fieldName string) {
	v.Check(NotBlank(password), fieldName, "must be provided")
	v.Check(MinChars(password, 8), fieldName, "must be at least 8 characters long")
	v.Check(MaxChars(password, 72), fieldName, "must not be more than 72 characters long")
}

// ValidateName validates name fields
func ValidateName(v *Validator, name, fieldName string) {
	v.Check(NotBlank(name), fieldName, "must be provided")
	v.Check(MinChars(name, 2), fieldName, "must be at least 2 characters long")
	v.Check(MaxChars(name, 50), fieldName, "must not be more than 50 characters long")
}

// ValidatePhone validates phone number format
func ValidatePhone(v *Validator, phone, fieldName string) {
	v.Check(NotBlank(phone), fieldName, "must be provided")
	v.Check(Matches(phone, PhoneRX), fieldName, "must be a valid phone number")
}

// ValidateRequired validates that a field is not empty
func ValidateRequired(v *Validator, value, fieldName string) {
	v.Check(NotBlank(value), fieldName, "must be provided")
}

// ValidateStringLength validates string length within bounds
func ValidateStringLength(v *Validator, value, fieldName string, min, max int) {
	if NotBlank(value) {
		v.Check(MinChars(value, min), fieldName, fmt.Sprintf("must be at least %d characters long", min))
		v.Check(MaxChars(value, max), fieldName, fmt.Sprintf("must not be more than %d characters long", max))
	}
}

// ValidateRole validates user role
func ValidateRole(v *Validator, role, fieldName string) {
	validRoles := []string{"admin", "manager", "staff", "guest"}
	v.Check(In(role, validRoles...), fieldName, "must be a valid role (admin, manager, staff, guest)")
}

// ValidateID validates that an ID is positive
func ValidateID(v *Validator, id int64, fieldName string) {
	v.Check(id > 0, fieldName, "must be a positive integer")
}

// ValidatePagination validates pagination parameters
func ValidatePagination(v *Validator, limit, offset int) {
	v.Check(limit > 0, "limit", "must be greater than 0")
	v.Check(limit <= 100, "limit", "must not be greater than 100")
	v.Check(offset >= 0, "offset", "must be greater than or equal to 0")
}
