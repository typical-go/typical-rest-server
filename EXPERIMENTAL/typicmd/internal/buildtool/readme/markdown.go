package readme

import (
	"fmt"
	"io"
)

// Markdown builder
type Markdown struct {
	io.Writer
}

// NewMarkdown return new instance of Markdown
func NewMarkdown(writer io.Writer) *Markdown {
	return &Markdown{
		Writer: writer,
	}
}

// OrderedList to create ordered list
func (m *Markdown) OrderedList(list ...string) {
	for i, item := range list {
		m.Writef("%d. %s\n", i+1, item)
	}
	m.WriteString("\n")
}

// InlineCode to create inline code
func (m *Markdown) InlineCode(code string) (int, error) {
	return m.WriteString("`" + code + "`")
}

// CodeSnippet to create code snippet
func (m *Markdown) CodeSnippet(code string) (int, error) {
	return m.WriteString("`" + code + "`")
}

func (m *Markdown) Writelnf(format string, args ...interface{}) (int, error) {
	return m.WriteString(fmt.Sprintf(format, args...) + "\n\n")
}

func (m *Markdown) Writeln(text string) (int, error) {
	return m.WriteString(text + "\n\n")
}

func (m *Markdown) Heading1(text string) (int, error) {
	return m.Writelnf("# %s", text)
}

func (m *Markdown) Heading2(text string) (int, error) {
	return m.Writelnf("## %s", text)
}

func (m *Markdown) Heading3(text string) (int, error) {
	return m.Writelnf("### %s", text)
}

func (m *Markdown) WriteString(s string) (int, error) {
	return m.Write([]byte(s))
}

func (m *Markdown) Comment(comment string) (int, error) {
	return m.Writelnf("<!-- %s -->", comment)
}

// Writef same WriteString with formatting
func (m *Markdown) Writef(format string, args ...interface{}) (int, error) {
	return m.WriteString(fmt.Sprintf(format, args...))
}
