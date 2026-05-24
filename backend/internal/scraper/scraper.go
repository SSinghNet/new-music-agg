package scraper

import (
	"context"
	"sync"

	"github.com/SSinghNet/new-music-agg/backend/internal/models"
)

// Run executes all scrapers concurrently across the given number of pages
// and streams results to the returned channel, which is closed when done.
func Run(ctx context.Context, pages int) <-chan *models.Release {
	out := make(chan *models.Release, 512)
	var wg sync.WaitGroup

	launch := func(fn func() []*models.Release) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range fn() {
				select {
				case out <- r:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	launch(func() []*models.Release { return GetAlbumsOfTheWeek(pages) })
	launch(func() []*models.Release { return GetAlbumsPremature(pages) })

	for i := range pages {
		page := i
		launch(func() []*models.Release { return GetAlbumOfTheDays(page) })
		launch(func() []*models.Release { return GetNewMusicFriday(page, 10) })
		launch(func() []*models.Release { return GetAlbums(page) })
		launch(func() []*models.Release { return GetTracks(page) })
		launch(func() []*models.Release { return GetBangers(page, 12) })
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
