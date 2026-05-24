package scraper

import (
	"testing"

	"github.com/SSinghNet/new-music-agg/internal/models"
)

func checkReleases(t *testing.T, releases []*models.Release) {
	t.Helper()
	if len(releases) == 0 {
		t.Fatal("expected releases, got none")
	}
	for _, r := range releases {
		if r.Name == "" {
			t.Errorf("release has empty name: %+v", r)
		}
		if len(r.Artists) == 0 {
			t.Errorf("release %q has no artists", r.Name)
		}
		if r.PublishDate.IsZero() {
			t.Errorf("release %q has zero date", r.Name)
		}
		if r.Link == "" {
			t.Errorf("release %q has no link", r.Name)
		}
	}
	t.Logf("scraped %d releases", len(releases))
}
