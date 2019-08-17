package gosrc

import (
	"io"
	"os"
	"sort"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe"
)

// SourceCode is source code recipe for generated.go in typical package
type SourceCode struct {
	PackageName  string
	Imports      []Import
	Structs      []Struct
	Constructors []string
	MockTargets  []string
	TestTargets  []string
}

func (r SourceCode) Write(w io.Writer) {
	writelnf(w, "// "+typirecipe.WaterMark+"\n")
	writelnf(w, "package %s", r.PackageName)

	for _, importPogo := range r.Imports {
		writelnf(w, `import %s "%s"`, importPogo.Alias, importPogo.PackageName)
	}

	for i := range r.Structs {
		r.Structs[i].Write(w)
	}

	writeln(w, "func init() {")
	for i := range r.Constructors {
		writelnf(w, "Context.AddConstructor(%s)", r.Constructors[i])
	}
	for i := range r.MockTargets {
		writelnf(w, "Context.AddMockTarget(\"%s\")", r.MockTargets[i])
	}
	for i := range r.TestTargets {
		writelnf(w, "Context.AddTestTarget(\"./%s\")", r.TestTargets[i])
	}
	writeln(w, "}")
}

// Cook to generate the recipe into file
func (r SourceCode) Cook(file string) (err error) {
	var f *os.File
	f, err = os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()

	r.sortOut()
	r.Write(f)

	return
}

// Blank is nothing to generate for recipe
func (r SourceCode) Blank() bool {
	return len(r.Imports) < 1 &&
		len(r.Structs) < 1 &&
		len(r.MockTargets) < 1 &&
		len(r.Constructors) < 1 &&
		len(r.TestTargets) < 1

}

func (r SourceCode) sortOut() {
	sort.Strings(r.Constructors)
	sort.Strings(r.MockTargets)
	sort.Strings(r.TestTargets)
}

// AddConstructorFunction to add FunctionPogo to constructor
func (r *SourceCode) AddConstructorFunction(pogos ...Function) {
	for _, pogo := range pogos {
		r.Constructors = append(r.Constructors, pogo.String())
	}
}

// AddConstructors to add constructors
func (r *SourceCode) AddConstructors(constructors ...string) {
	r.Constructors = append(r.Constructors, constructors...)
}

// AddMockTargets to add constructors
func (r *SourceCode) AddMockTargets(mockTargets ...string) {
	r.MockTargets = append(r.MockTargets, mockTargets...)
}

// AddTestTargets to add constructors
func (r *SourceCode) AddTestTargets(testTargets ...string) {
	r.TestTargets = append(r.TestTargets, testTargets...)
}

// AddImport to add import POGO
func (r *SourceCode) AddImport(pogos ...Import) {
	r.Imports = append(r.Imports, pogos...)
}
