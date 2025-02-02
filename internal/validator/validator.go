package validator

import (
	"regexp"
)

// Declare a regular expression for sanity checking the format of email addresses.
var (
	EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
)

// Validator type contains a map of validation errors.
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance with an empty errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for the given key).
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

// PermittedValue returns true if a specific value is in a list.
func PermittedValue(value interface{}, permittedValues ...interface{}) bool {
	for _, v := range permittedValues {
		if value == v {
			return true
		}
	}
	return false
}

// Unique returns true if all values in a slice are unique.
func Unique(values []string) bool {
	uniqueValues := make(map[interface{}]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
