package store

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/SSinghNet/new-music-agg/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GORMStore struct {
	db *gorm.DB
}

// resolveIPv4DSN resolves the hostname in a postgres DSN to an IPv4 address.
// Supabase direct connections fail over IPv6 on some networks.
func resolveIPv4DSN(dsn string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn, nil // not a URL-style DSN; leave unchanged
	}
	host := u.Hostname()
	if host == "" || net.ParseIP(host) != nil {
		return dsn, nil // already an IP; skip lookup
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", fmt.Errorf("lookup %s: %w", host, err)
	}
	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip != nil && ip.To4() != nil {
			port := u.Port()
			if port != "" {
				u.Host = ip.String() + ":" + port
			} else {
				u.Host = ip.String()
			}
			return u.String(), nil
		}
	}
	return dsn, nil // no IPv4 found; fall through to original
}

func NewPostgresStore(dsn string) (*GORMStore, error) {
	ipv4DSN, err := resolveIPv4DSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("resolveIPv4DSN: %w", err)
	}

	cfg, err := pgx.ParseConfig(ipv4DSN)
	if err != nil {
		return nil, fmt.Errorf("pgx.ParseConfig: %w", err)
	}
	// Avoid "prepared statement already exists" when database/sql reuses connections
	// across goroutines — simple protocol skips server-side prepared statement caching.
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	sqlDB := stdlib.OpenDB(*cfg)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	db, err := gorm.Open(gorm_postgres.New(gorm_postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}
	return &GORMStore{db: db}, nil
}

func (s *GORMStore) Close() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

var releaseConflict = clause.OnConflict{
	Columns: []clause.Column{{Name: "link"}},
	TargetWhere: clause.Where{Exprs: []clause.Expression{
		gorm.Expr("link <> ''"),
	}},
	DoUpdates: clause.AssignmentColumns([]string{
		"name", "publish_date", "source", "release_type", "special", "updated_at",
	}),
}

// findOrCreateArtists resolves artist names to DB records in two round-trips:
// one bulk SELECT for existing rows, one batch INSERT (ON CONFLICT DO NOTHING) for new ones.
func (s *GORMStore) findOrCreateArtists(ctx context.Context, artists []models.Artist) ([]models.Artist, error) {
	if len(artists) == 0 {
		return nil, nil
	}

	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}

	var existing []models.Artist
	if err := s.db.WithContext(ctx).Where("name IN ?", names).Find(&existing).Error; err != nil {
		return nil, err
	}

	byName := make(map[string]models.Artist, len(existing))
	for _, a := range existing {
		byName[a.Name] = a
	}

	var toCreate []models.Artist
	for _, name := range names {
		if _, ok := byName[name]; !ok {
			toCreate = append(toCreate, models.Artist{Name: name})
		}
	}

	if len(toCreate) > 0 {
		// ON CONFLICT DO NOTHING handles concurrent inserts of the same artist name.
		if err := s.db.WithContext(ctx).
			Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}}, DoNothing: true}).
			Create(&toCreate).Error; err != nil {
			return nil, err
		}

		// Re-fetch any rows that were skipped by DO NOTHING (ID will be 0).
		var missing []string
		for _, a := range toCreate {
			if a.ID == 0 {
				missing = append(missing, a.Name)
			} else {
				byName[a.Name] = a
			}
		}
		if len(missing) > 0 {
			var refetched []models.Artist
			if err := s.db.WithContext(ctx).Where("name IN ?", missing).Find(&refetched).Error; err != nil {
				return nil, err
			}
			for _, a := range refetched {
				byName[a.Name] = a
			}
		}
	}

	resolved := make([]models.Artist, 0, len(names))
	for _, name := range names {
		resolved = append(resolved, byName[name])
	}
	return resolved, nil
}

func (s *GORMStore) Upsert(ctx context.Context, r *models.Release) error {
	artists, err := s.findOrCreateArtists(ctx, r.Artists)
	if err != nil {
		return fmt.Errorf("find/create artists: %w", err)
	}

	// Upsert the release row without associations
	saved := *r
	saved.Artists = nil
	if err := s.db.WithContext(ctx).Clauses(releaseConflict).Create(&saved).Error; err != nil {
		return err
	}
	r.ID = saved.ID

	// Replace artist associations for this release
	return s.db.WithContext(ctx).Model(r).Association("Artists").Replace(artists)
}

func (s *GORMStore) UpsertBatch(ctx context.Context, releases []*models.Release) error {
	for _, r := range releases {
		if err := s.Upsert(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (s *GORMStore) List(ctx context.Context, p ListParams) ([]*models.Release, int, error) {
	q := s.db.WithContext(ctx).Model(&models.Release{})

	if p.Source != nil {
		q = q.Where("releases.source = ?", *p.Source)
	}
	if p.ReleaseType != nil {
		q = q.Where("releases.release_type = ?", *p.ReleaseType)
	}
	if p.Special != nil {
		q = q.Where("releases.special = ?", *p.Special)
	}
	if p.Artist != nil {
		q = q.Joins("JOIN release_artists ra ON ra.release_id = releases.id").
			Joins("JOIN artists a ON a.id = ra.artist_id").
			Where("a.name = ?", *p.Artist)
	}
	if p.DateFrom != nil {
		q = q.Where("releases.publish_date >= ?", *p.DateFrom)
	}
	if p.DateTo != nil {
		q = q.Where("releases.publish_date <= ?", *p.DateTo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dir := "DESC"
	if p.OrderDir == "asc" {
		dir = "ASC"
	}
	limit := p.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var releases []*models.Release
	err := q.Preload("Artists").
		Order("releases.publish_date " + dir).
		Limit(limit).Offset(p.Offset).
		Find(&releases).Error
	return releases, int(total), err
}

func (s *GORMStore) GetByID(ctx context.Context, id uint) (*models.Release, error) {
	var r models.Release
	err := s.db.WithContext(ctx).Preload("Artists").First(&r, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &r, err
}
