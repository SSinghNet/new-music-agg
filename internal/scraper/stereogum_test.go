package scraper

import "testing"

func TestGetAlbumsOfTheWeek(t *testing.T) {
	checkReleases(t, GetAlbumsOfTheWeek(5))
}

func TestGetAlbumsPremature(t *testing.T) {
	checkReleases(t, GetAlbumsPremature(5))
}
