package markdown

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

// UnorderedList to create unorder list
func (m *Markdown) UnorderedList(list ...string) {
	for _, item := range list {
		m.Writef("- %s\n", item)
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

// Writelnf write formatted string followed by end line
func (m *Markdown) Writelnf(format string, args ...interface{}) (int, error) {
	return m.WriteString(fmt.Sprintf(format, args...) + "\n\n")
}

// Writeln write string followed by end line
func (m *Markdown) Writeln(text string) (int, error) {
	return m.WriteString(text + "\n\n")
}

// H1 write for heading 1
func (m *Markdown) H1(text string) (int, error) {
	return m.Writelnf("# %s", text)
}

// H2 write for heading 2
func (m *Markdown) H2(text string) (int, error) {
	return m.Writelnf("## %s", text)
}

// H3 write for heading 3
func (m *Markdown) H3(text string) (int, error) {
	return m.Writelnf("### %s", text)
}

// H4 write for heading 4
func (m *Markdown) H4(text string) (int, error) {
	return m.Writelnf("#### %s", text)
}

// WriteString to write string
func (m *Markdown) WriteString(s string) (int, error) {
	return m.Write([]byte(s))
}

// Comment for comment
func (m *Markdown) Comment(comment string) (int, error) {
	return m.Writelnf("<!-- %s -->", comment)
}

// Writef same WriteString with formatting
func (m *Markdown) Writef(format string, args ...interface{}) (int, error) {
	return m.WriteString(fmt.Sprintf(format, args...))
}
