package typicli

import (
	"io/ioutil"
	"strings"

	"github.com/typical-go/typical-rest-server/experimental/typienv"
)

func (t *TypicalCli) bundleCliSideEffects() (err error) {

	builder := &strings.Builder{}
	builder.WriteString("package main\n")
	builder.WriteString("import(\n")

	for _, module := range t.Modules {
		for _, sideEffect := range module.SideEffects {
			builder.WriteString("_ \"" + sideEffect + "\"\n")
		}

		for _, sideEffect := range module.TypiCliSideEffects {
			builder.WriteString("_ \"" + sideEffect + "\"\n")
		}
	}
	builder.WriteString(")")

	filename := typienv.TypicalMainPackage() + "/init.go"
	err = ioutil.WriteFile(filename, []byte(builder.String()), 0644)
	if err != nil {
		return
	}

	runOrFatalSilently(goCommand(), "fmt", filename)

	return

}
