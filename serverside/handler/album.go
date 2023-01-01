package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/miko2823/go-docker/pkg"
	"github.com/miko2823/go-docker/usecase"
)

type AlbumHandler interface {
	RegisterHandlers() http.Handler
	getAlbum(w http.ResponseWriter, r *http.Request)
}

type albumHandler struct {
	albumUsecase usecase.AlbumUsecase
}

func NewAlbumHandler(usecase usecase.AlbumUsecase) AlbumHandler {
	return &albumHandler{usecase}
}

func (p albumHandler) getAlbum(w http.ResponseWriter, r *http.Request) {

	// var requestPayload struct {
	// 	Id string `json:"id"`
	// }
	// err := pkg.ReadJSON(w, r, &requestPayload)
	// fmt.Println(err)
	// if err != nil {
	// 	pkg.ErrorJSON(w, err, http.StatusBadRequest)
	// }

	id := strings.TrimPrefix(r.URL.Path, "/v1/album/")
	album, err := p.albumUsecase.Get(id)

	if err != nil {
		pkg.ErrorJSON(w, err, http.StatusBadRequest)
	}

	payload := pkg.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("get album data"),
		Data:    album,
	}
	pkg.WriteJson(w, http.StatusAccepted, payload)
}
