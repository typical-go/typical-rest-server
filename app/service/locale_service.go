package service

import (
	"github.com/typical-go/typical-rest-server/app/repository"
	"go.uber.org/dig"
)

// LocaleService contain logic for LocaleController
type LocaleService interface {
	repository.LocaleRepo
}

// LocaleServiceImpl is implementation of LocaleService
type LocaleServiceImpl struct {
	dig.In
	repository.LocaleRepo
}

// NewLocaleService return new instance of LocaleService
func NewLocaleService(impl LocaleServiceImpl) LocaleService {
	return &impl
}
