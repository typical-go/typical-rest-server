package typrest

const serviceTmpl = `package service

import (
	"github.com/typical-go/typical-rest-server/app/repository"
	"go.uber.org/dig"
)

// {{ .Name }}Service contain logic for {{ .Name }} Controller
type {{ .Name }}Service interface {
	repository.{{ .Name }}Repo
}

// {{ .Name }}ServiceImpl is implementation of {{ .Name }}Service
type {{ .Name }}ServiceImpl struct {
	dig.In
	repository.{{ .Name }}Repo
}

// New{{ .Name }}Service return new instance of {{ .Name }}Service
func New{{ .Name }}Service(impl {{ .Name }}ServiceImpl) {{ .Name }}Service {
	return &impl
}
`
