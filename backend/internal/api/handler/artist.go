package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/SSinghNet/new-music-agg/backend/internal/httputil"
	"github.com/SSinghNet/new-music-agg/backend/internal/models"
	"github.com/SSinghNet/new-music-agg/backend/internal/service"
	"github.com/SSinghNet/new-music-agg/backend/internal/store"

	"github.com/go-chi/chi/v5"
)

type ArtistHandler struct {
	svc service.ArtistService
}

func NewArtistHandler(svc service.ArtistService) *ArtistHandler {
	return &ArtistHandler{svc: svc}
}

func (h *ArtistHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p := store.ListArtistsParams{}

	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.Limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.Offset = n
		}
	}

	artists, total, err := h.svc.List(r.Context(), p)
	if err != nil {
		log.Printf("list artists: %v", err)
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if artists == nil {
		artists = []*models.Artist{}
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 50
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"data": artists,
		"meta": httputil.Meta{Total: total, Limit: limit, Offset: p.Offset},
	})
}

func (h *ArtistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	artist, err := h.svc.GetByID(r.Context(), uint(id))
	if errors.Is(err, store.ErrArtistNotFound) {
		httputil.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if err != nil {
		log.Printf("get artist %d: %v", id, err)
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, artist)
}
