package scraper

import "testing"

func TestGetBangers(t *testing.T) {
	all := GetBangers(1, 12)
	for page := 2; page <= 5; page++ {
		all = append(all, GetBangers(page, 12)...)
	}
	checkReleases(t, all)
}
