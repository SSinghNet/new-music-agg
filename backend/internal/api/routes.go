package api

import (
	"net/http"

	"github.com/SSinghNet/new-music-agg/internal/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(releases *handler.ReleaseHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/releases", releases.List)
		r.Get("/releases/{id}", releases.GetByID)
	})

	return r
}
