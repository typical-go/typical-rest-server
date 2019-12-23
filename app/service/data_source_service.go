package service

import (
	"github.com/typical-go/typical-rest-server/app/repository"
	"go.uber.org/dig"
)

// DataSourceService contain logic for DataSourceController
type DataSourceService interface {
	repository.DataSourceRepo
}

// DataSourceServiceImpl is implementation of DataSourceService
type DataSourceServiceImpl struct {
	dig.In
	repository.DataSourceRepo
}

// NewDataSourceService return new instance of DataSourceService
func NewDataSourceService(impl DataSourceServiceImpl) DataSourceService {
	return &impl
}
