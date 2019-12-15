package typrest

const serviceTmpl = `package service

import (
	"{{.ProjectPackage}}/app/repository"
	"go.uber.org/dig"
)

// {{.TypeName}}Service contain logic for {{.TypeName}}Controller
type {{.TypeName}}Service interface {
	repository.{{.TypeName}}Repo
}

// {{.TypeName}}ServiceImpl is implementation of {{.TypeName}}Service
type {{.TypeName}}ServiceImpl struct {
	dig.In
	repository.{{.TypeName}}Repo
}

// New{{.TypeName}}Service return new instance of {{.TypeName}}Service
func New{{.TypeName}}Service(impl {{.TypeName}}ServiceImpl) {{.TypeName}}Service {
	return &impl
}
`
