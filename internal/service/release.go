package service

import (
	"context"

	"github.com/SSinghNet/new-music-agg/internal/models"
	"github.com/SSinghNet/new-music-agg/internal/store"
)

type ReleaseService interface {
	List(ctx context.Context, p store.ListParams) ([]*models.Release, int, error)
	GetByID(ctx context.Context, id uint) (*models.Release, error)
}

type releaseService struct {
	store store.Store
}

func NewReleaseService(st store.Store) ReleaseService {
	return &releaseService{store: st}
}

func (s *releaseService) List(ctx context.Context, p store.ListParams) ([]*models.Release, int, error) {
	return s.store.List(ctx, p)
}

func (s *releaseService) GetByID(ctx context.Context, id uint) (*models.Release, error) {
	return s.store.GetByID(ctx, id)
}
