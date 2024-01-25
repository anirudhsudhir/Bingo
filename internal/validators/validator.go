package validators

import (
	"unicode/utf8"
)

type Validator struct {
	FormErrors map[string]string
}

func (v *Validator) ValidateElement(valid bool, key, message string) {
	if v.FormErrors == nil {
		v.FormErrors = make(map[string]string)
	}
	if !valid {
		_, present := v.FormErrors[key]
		if !present {
			v.FormErrors[key] = message
		}
	}
}

func (v *Validator) ValidForm() bool {
	return len(v.FormErrors) == 0
}

func NoContent(element string) bool {
	return len(element) > 0
}

func MaxLen(element string, maxSize int) bool {
	return utf8.RuneCountInString(element) <= maxSize
}

func AllowedValues(element int, allowedVals ...int) bool {
	for i := range allowedVals {
		if allowedVals[i] == element {
			return true
		}
	}
	return false
}
