package generated

import (
	"fmt"
	"strings"
)

type Builder struct {
	strings.Builder
}

func (b *Builder) Printlnf(format string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(format+"\n", args...))
}

func (b *Builder) Printf(format string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(format, args...))
}
