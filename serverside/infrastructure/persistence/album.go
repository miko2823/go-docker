package persistence

import (
	"database/sql"
	"time"

	"github.com/miko2823/go-docker/domain/album/models"
	"github.com/miko2823/go-docker/domain/album/repository"
)

type albumPersistence struct {
	Conn *sql.DB
}

func NewAlbumPersistence(conn *sql.DB) repository.AlbumRepository {
	return albumPersistence{conn}
}

func (r albumPersistence) Get(id string) (models.Album, error) {
	// TODO connenct repo
	// album := models.Album{ID: "123", Name: "album1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	row := r.Conn.QueryRow("SELECT * FROM users WHERE id = %s", id)
	return convertToAlbum(row)
}

func convertToAlbum(row *sql.Row) (models.Album, error) {
	album := models.Album{ID: "123", Name: "Tokyo visit", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	// err := row.Scan(&album.ID, &album.Name, &album.CreatedAt, &album.UpdatedAt)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }
	return album, nil
}
