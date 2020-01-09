package service

import (
	"github.com/typical-go/typical-rest-server/app/repository"
	"go.uber.org/dig"
)

// MusicService contain logic for MusicController [mock]
type MusicService interface {
	repository.MusicRepo
}

// MusicServiceImpl is implementation of MusicService
type MusicServiceImpl struct {
	dig.In
	repository.MusicRepo
}

// NewMusicService return new instance of MusicService [autowire]
func NewMusicService(impl MusicServiceImpl) MusicService {
	return &impl
}
