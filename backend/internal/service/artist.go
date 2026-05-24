package service

import (
	"context"

	"github.com/SSinghNet/new-music-agg/backend/internal/models"
	"github.com/SSinghNet/new-music-agg/backend/internal/store"
)

type ArtistService interface {
	List(ctx context.Context, p store.ListArtistsParams) ([]*models.Artist, int, error)
	GetByID(ctx context.Context, id uint) (*models.Artist, error)
}

type artistService struct {
	store store.Store
}

func NewArtistService(st store.Store) ArtistService {
	return &artistService{store: st}
}

func (s *artistService) List(ctx context.Context, p store.ListArtistsParams) ([]*models.Artist, int, error) {
	return s.store.ListArtists(ctx, p)
}

func (s *artistService) GetByID(ctx context.Context, id uint) (*models.Artist, error) {
	return s.store.GetArtistByID(ctx, id)
}
