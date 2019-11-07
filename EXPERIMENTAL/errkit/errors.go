package errkit

import "strings"

// Errors to handle multiple error
type Errors []error

// Add the error if not nil
func (e *Errors) Add(err error) bool {
	if err == nil {
		return false
	}
	*e = append(*e, err)
	return true
}

func (e Errors) Error() string {
	var builder strings.Builder
	for i, err := range e {
		if i > 0 {
			builder.WriteString("; ")
		}
		builder.WriteString(err.Error())
	}
	return builder.String()
}

// IsEmpty return true if no error
func (e Errors) IsEmpty() bool {
	return len(e) < 1
}
