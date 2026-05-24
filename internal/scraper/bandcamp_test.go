package scraper

import "testing"

func TestGetAlbumOfTheDays(t *testing.T) {
	all := GetAlbumOfTheDays(1)
	for page := 2; page <= 5; page++ {
		all = append(all, GetAlbumOfTheDays(page)...)
	}
	checkReleases(t, all)
}
