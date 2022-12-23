package persistence

import (
	"time"

	"github.com/miko2823/go-docker/domain/album/models"
	"github.com/miko2823/go-docker/domain/album/repository"
)

type albumPersistence struct{}

func NewAlbumPersistence() repository.AlbumRepository {
	return albumPersistence{}
}

func (r albumPersistence) Get(id string) (models.Album, error) {
	// TODO connenct repo
	album := models.Album{ID: "123", Name: "album1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return album, nil
}
