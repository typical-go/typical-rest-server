package service

import (
	"github.com/typical-go/typical-rest-server/app/repository"
	"go.uber.org/dig"
)

// TagService contain logic for TagController
type TagService interface {
	repository.TagRepo
}

// TagServiceImpl is implementation of TagService
type TagServiceImpl struct {
	dig.In
	repository.TagRepo
}

// NewTagService return new instance of TagService
func NewTagService(impl TagServiceImpl) TagService {
	return &impl
}
