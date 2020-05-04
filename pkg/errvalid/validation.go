package errvalid

import (
	"fmt"
	"strings"
)

const (
	validationErrPrefix = "Validation: "
	validationSep       = "\n"
)

// New validation error from messages
func New(msgs ...string) error {
	return fmt.Errorf("%s%s", validationErrPrefix, strings.Join(msgs, validationSep))
}

// Wrap error to validation error
func Wrap(err error) error {
	return fmt.Errorf("%s%w", validationErrPrefix, err)
}

// Messages validation
func Messages(err error) []string {
	msg := Message(err)
	if msg == "" {
		return []string{}
	}

	return strings.Split(msg, validationSep)
}

// Message of validation
func Message(err error) string {
	if !Check(err) {
		return ""
	}

	msg := err.Error()
	return strings.TrimSpace(msg[len(validationErrPrefix):])
}

// Check return true if err is validation error
func Check(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), validationErrPrefix)
}
