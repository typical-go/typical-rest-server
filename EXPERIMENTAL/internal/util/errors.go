package util

import "strings"

// Errors to handle multiple error
type Errors []error

// Add the error if not nil
func (e Errors) Add(err error) bool {
	if err == nil {
		return false
	}

	e = append(e, err)
	return true
}

func (e Errors) Error() string {
	var builder strings.Builder
	for i := range e {
		builder.WriteString(e[i].Error())
		builder.WriteString("; ")
	}

	return builder.String()
}
