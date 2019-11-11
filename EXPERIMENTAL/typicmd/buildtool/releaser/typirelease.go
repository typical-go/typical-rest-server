package releaser

import (
	"strings"
)

func clean(message string) string {
	iCoAuthor := strings.Index(message, "Co-Authored-By")
	if iCoAuthor > 0 {
		message = message[0:strings.Index(message, "Co-Authored-By")]
	}
	message = strings.TrimSpace(message)
	return message
}

func ignoring(changelog string) bool {
	message := clean(message(changelog))
	lowerMessage := strings.ToLower(message)
	return strings.HasPrefix(lowerMessage, "merge") ||
		strings.HasPrefix(lowerMessage, "bump") ||
		strings.HasPrefix(lowerMessage, "revision") ||
		strings.HasPrefix(lowerMessage, "generate") ||
		strings.HasPrefix(lowerMessage, "wip")
}

func message(changelog string) string {
	if len(changelog) < 7 {
		return ""
	}
	return strings.TrimSpace(changelog[7:])
}
