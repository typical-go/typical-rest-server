package typimain

import (
	"io/ioutil"
	"strings"

	"github.com/typical-go/typical-rest-server/experimental/typienv"
)

const (
	mainInitFile = "init.go"
)

func (t *TypicalTask) bundleCliSideEffects() error {
	var sideEffects []string
	for _, module := range t.Modules {
		sideEffects = append(sideEffects, module.SideEffects...)
		sideEffects = append(sideEffects, module.TypiCliSideEffects...)
	}
	filename := typienv.TypicalMainPackage() + "/" + mainInitFile
	return bundleSideEffects(filename, sideEffects)
}

func (t *TypicalTask) bundleAppSideEffects() error {
	var sideEffects []string
	for _, module := range t.Modules {
		sideEffects = append(sideEffects, module.SideEffects...)
		sideEffects = append(sideEffects, module.TypiAppSideEffects...)
	}

	filename := typienv.MainPackage(t.AppPkgOrDefault()) + "/" + mainInitFile
	return bundleSideEffects(filename, sideEffects)
}

func bundleSideEffects(filename string, sideEffects []string) (err error) {

	builder := &strings.Builder{}
	builder.WriteString("package main\n")
	builder.WriteString("import(\n")

	for _, sideEffect := range sideEffects {
		builder.WriteString("_ \"" + sideEffect + "\"\n")
	}

	builder.WriteString(")")

	err = ioutil.WriteFile(filename, []byte(builder.String()), 0644)
	if err != nil {
		return
	}

	runOrFatalSilently(goCommand(), "fmt", filename)

	return

}
