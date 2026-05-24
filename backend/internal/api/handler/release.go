package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SSinghNet/new-music-agg/internal/httputil"
	"github.com/SSinghNet/new-music-agg/internal/models"
	"github.com/SSinghNet/new-music-agg/internal/service"
	"github.com/SSinghNet/new-music-agg/internal/store"

	"github.com/go-chi/chi/v5"
)

type ReleaseHandler struct {
	svc service.ReleaseService
}

func NewReleaseHandler(svc service.ReleaseService) *ReleaseHandler {
	return &ReleaseHandler{svc: svc}
}

func (h *ReleaseHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p := store.ListParams{
		OrderDir: q.Get("order_dir"),
	}

	if v := q.Get("source"); v != "" {
		s := models.SourceName(v)
		p.Source = &s
	}
	if v := q.Get("type"); v != "" {
		t := models.ReleaseType(v)
		p.ReleaseType = &t
	}
	if v := q.Get("special"); v != "" {
		sl := models.SpecialLabel(v)
		p.Special = &sl
	}
	if v := q.Get("artist"); v != "" {
		p.Artist = &v
	}
	if v := q.Get("from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid from date, use YYYY-MM-DD"})
			return
		}
		p.DateFrom = &t
	}
	if v := q.Get("to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid to date, use YYYY-MM-DD"})
			return
		}
		p.DateTo = &t
	}
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

	releases, total, err := h.svc.List(r.Context(), p)
	if err != nil {
		log.Printf("list releases: %v", err)
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if releases == nil {
		releases = []*models.Release{}
	}

	limit := p.Limit
	if limit <= 0 {
		limit = 50
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"data": releases,
		"meta": httputil.Meta{Total: total, Limit: limit, Offset: p.Offset},
	})
}

func (h *ReleaseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	release, err := h.svc.GetByID(r.Context(), uint(id))
	if errors.Is(err, store.ErrNotFound) {
		httputil.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if err != nil {
		log.Printf("get release %d: %v", id, err)
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, release)
}
