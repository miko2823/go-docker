package repository

import (
	"github.com/miko2823/go-docker/domain/album/models"
)

type AlbumRepository interface {
	Get(id string) (models.Album, error)
}
