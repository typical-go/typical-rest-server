package service

import (
	"github.com/typical-go/typical-rest-server/app/repository"
	"go.uber.org/dig"
)

// MusicService contain logic for Music Controller
type MusicService interface {
	repository.MusicRepo
}

// MusicServiceImpl is implementation of MusicService
type MusicServiceImpl struct {
	dig.In
	repository.MusicRepo
}

// NewMusicService return new instance of MusicService
func NewMusicService(impl MusicServiceImpl) MusicService {
	return &impl
}
