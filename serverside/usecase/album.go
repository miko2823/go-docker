package usecase

import (
	"github.com/miko2823/go-docker/domain/album/models"
	"github.com/miko2823/go-docker/domain/album/repository"
)

type AlbumUsecase interface {
	Get(id string) (models.Album, error)
}

type albumUsecase struct {
	albumRepository repository.AlbumRepository
}

func NewAlbumUsecase(repo repository.AlbumRepository) AlbumUsecase {
	return albumUsecase{repo}
}

func (u albumUsecase) Get(id string) (models.Album, error) {
	album, err := u.albumRepository.Get(id)
	if err != nil {
		return models.Album{}, err
	}
	return album, nil
}
