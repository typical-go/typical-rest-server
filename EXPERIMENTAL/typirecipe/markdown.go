package typirecipe

import (
	"fmt"
	"strings"
)

// Markdown builder
type Markdown struct {
	strings.Builder
}

// OrderedList to create ordered list
func (m *Markdown) OrderedList(list ...string) {
	for i, item := range list {
		m.Writef("%d. %s\n", i, item)
	}
}

// InlineCode to create inline code
func (m *Markdown) InlineCode(code string) (int, error) {
	return m.WriteString("`" + code + "`")
}

// CodeSnippet to create code snippet
func (m *Markdown) CodeSnippet(code string) (int, error) {
	return m.WriteString("`" + code + "`")
}

// Writef same WriteString with formatting
func (m *Markdown) Writef(format string, args ...interface{}) (int, error) {
	return m.WriteString(fmt.Sprintf(format, args...))
}

// Writefln same with WriteString with formatting and new line
func (m *Markdown) Writefln(format string, args ...interface{}) (int, error) {
	return m.WriteString(fmt.Sprintf(format, args...) + "\n")
}

// Writeln same with WriteString with new line
func (m *Markdown) Writeln(text string) (int, error) {
	return m.WriteString(text + "\n")
}
