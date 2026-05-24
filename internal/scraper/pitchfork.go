package scraper

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SSinghNet/new-music-agg/internal/models"

	"github.com/gocolly/colly"
)

func makePitchforkRelease(name string, artistNames []string, publishDate string, link string, special string, releaseType models.ReleaseType) *models.Release {
	release := &models.Release{}

	name = strings.ReplaceAll(name, "“", "")
	name = strings.ReplaceAll(name, "”", "")
	name = strings.ReplaceAll(name, "’", "'")
	release.Name = strings.TrimSpace(name)

	for _, n := range artistNames {
		n = strings.TrimSpace(n)
		if n != "" {
			release.Artists = append(release.Artists, models.Artist{Name: n})
		}
	}
	if len(release.Artists) == 0 {
		return nil
	}

	date, err := time.Parse("January 2, 2006", publishDate)
	if err != nil {
		log.Printf("pitchfork: parse date %q: %v", publishDate, err)
		return nil
	}
	release.PublishDate = date
	release.Link = "https://pitchfork.com" + link

	special = strings.TrimSpace(special)
	if special != "" {
		s := models.SpecialLabel(special)
		release.Special = &s
	}

	release.ReleaseType = releaseType
	release.Source = models.SourcePitchfork

	return release
}

func GetAlbums(page int) []*models.Release {
	c := colly.NewCollector()
	var albums []*models.Release
	c.OnHTML(".summary-item", func(e *colly.HTMLElement) {
		name := e.ChildText(".summary-item__hed")
		artists := strings.Split(e.ChildText(".summary-item__sub-hed"), "/")
		publishDate := e.ChildText(".summary-item__publish-date")
		link := e.ChildAttr(".summary-item__hed-link", "href")
		best := e.ChildText("[class^='SummaryItemReviewLabelWrapper-']")

		r := makePitchforkRelease(name, artists, publishDate, link, best, models.TypeAlbum)
		if r != nil {
			albums = append(albums, r)
		}
	})
	c.Visit("https://pitchfork.com/reviews/albums/?page=" + strconv.Itoa(page))

	return albums
}

func GetTracks(page int) []*models.Release {
	c := colly.NewCollector()
	var tracks []*models.Release
	c.OnHTML(".summary-item", func(e *colly.HTMLElement) {
		name := e.ChildText(".summary-item__hed")
		artists := strings.Split(e.ChildText(".summary-item__sub-hed"), "/")
		publishDate := e.ChildText(".summary-item__publish-date")
		link := e.ChildAttr(".summary-item__hed-link", "href")
		best := e.ChildText("[class^='SummaryItemReviewLabelWrapper-']")

		r := makePitchforkRelease(name, artists, publishDate, link, best, models.TypeTrack)
		if r != nil {
			tracks = append(tracks, r)
		}
	})
	c.Visit("https://pitchfork.com/reviews/tracks/?page=" + strconv.Itoa(page))

	return tracks
}
