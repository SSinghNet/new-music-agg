package scraper

import "testing"

func TestGetAlbums(t *testing.T) {
	all := GetAlbums(1)
	for page := 2; page <= 5; page++ {
		all = append(all, GetAlbums(page)...)
	}
	checkReleases(t, all)
}

func TestGetTracks(t *testing.T) {
	all := GetTracks(1)
	for page := 2; page <= 5; page++ {
		all = append(all, GetTracks(page)...)
	}
	checkReleases(t, all)
}
