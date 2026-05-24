package store

import (
	"context"
	"errors"
	"time"

	"github.com/SSinghNet/new-music-agg/backend/internal/models"
)

type ListParams struct {
	Source      *models.SourceName
	ReleaseType *models.ReleaseType
	Special     *models.SpecialLabel
	Artist      *string
	DateFrom    *time.Time
	DateTo      *time.Time
	Limit       int
	Offset      int
	OrderDir    string // "desc" (default) | "asc"
}

type ListArtistsParams struct {
	Limit  int
	Offset int
}

type Store interface {
	Upsert(ctx context.Context, r *models.Release) error
	UpsertBatch(ctx context.Context, releases []*models.Release) error
	List(ctx context.Context, p ListParams) ([]*models.Release, int, error)
	GetByID(ctx context.Context, id uint) (*models.Release, error)
	ListArtists(ctx context.Context, p ListArtistsParams) ([]*models.Artist, int, error)
	GetArtistByID(ctx context.Context, id uint) (*models.Artist, error)
	Close()
}

var ErrNotFound = errors.New("release not found")
var ErrArtistNotFound = errors.New("artist not found")
