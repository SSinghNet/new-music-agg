package scraper

import "testing"

func TestGetNewMusicFriday(t *testing.T) {
	all := GetNewMusicFriday(0, 5)
	for page := 1; page <= 4; page++ {
		all = append(all, GetNewMusicFriday(page, 5)...)
	}
	checkReleases(t, all)
}
