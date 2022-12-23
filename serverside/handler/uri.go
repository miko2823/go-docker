package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (p albumHandler) RegisterHandlers() http.Handler {
	fmt.Println("RegisterHandlers")
	r := chi.NewRouter()
	r.Get("/{id}", p.getAlbum)
	return r
}
