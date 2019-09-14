package golang

import (
	"fmt"
	"io"
)

// Initialization to generate init() function
type Initialization struct {
	Statements []string
}

// IsBlank return true if no statement
func (i *Initialization) IsBlank() bool {
	return len(i.Statements) < 1
}

// AddStatement to add statement to initialization
func (i *Initialization) AddStatement(stmt string) {
	i.Statements = append(i.Statements, stmt)
}

// AddConstructors to add constructors
func (i *Initialization) AddConstructors(constructors ...string) {
	for _, constructor := range constructors {
		i.AddStatement(fmt.Sprintf("typical.Context.Constructors.Add(%s)", constructor))
	}
}

// AddMockTargets to add constructors
func (i *Initialization) AddMockTargets(mockTargets ...string) {
	for _, mockTarget := range mockTargets {
		i.AddStatement(fmt.Sprintf("typical.Context.MockTargets.Add(\"%s\")", mockTarget))
	}
}

// AddTestTargets to add constructors
func (i *Initialization) AddTestTargets(testTargets ...string) {
	for _, testTarget := range testTargets {
		i.AddStatement(fmt.Sprintf("typical.Context.TestTargets.Add(\"./%s\")", testTarget))
	}
}

func (i Initialization) Write(w io.Writer) {
	writeln(w, "func init() {")
	for _, stmt := range i.Statements {
		writeln(w, stmt)
	}
	writeln(w, "}")
}
