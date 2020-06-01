package tmpl

// Service temmplate
const Service = `package service

import (
	"{{.ProjectPackage}}/app/repository"
	"go.uber.org/dig"
)

// {{.Type}}Service contain logic for {{.Type}}Controller
type {{.Type}}Service interface {
	repository.{{.Type}}Repo
}

// {{.Type}}ServiceImpl is implementation of {{.Type}}Service
type {{.Type}}ServiceImpl struct {
	dig.In
	repository.{{.Type}}Repo
}

// New{{.Type}}Service return new instance of {{.Type}}Service
func New{{.Type}}Service(impl {{.Type}}ServiceImpl) {{.Type}}Service {
	return &impl
}
`
