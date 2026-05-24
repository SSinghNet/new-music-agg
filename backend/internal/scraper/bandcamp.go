package scraper

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/SSinghNet/new-music-agg/backend/internal/models"

	"github.com/gocolly/colly"
)

// After normalization all quotes are straight " and all whitespace is ASCII space.
var bandcampTitleRe = regexp.MustCompile(`^(.+?)\s*[,，]\s*"([^"]+)"`)
var bandcampArtistSplitRe = regexp.MustCompile(`\s*(?:,|&|\band\b|x)\s*`)

// normalizeTitle collapses Unicode whitespace variants to ASCII space and
// maps curly/fancy quotes to straight double quotes so the regex stays simple.
func normalizeTitle(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case unicode.IsSpace(r):
			b.WriteByte(' ')
		case r == '“' || r == '‘' || r == '’':
			b.WriteByte('"')
		case r == '”':
			b.WriteByte('"')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func makeBandcampRelease(name string, artistNames []string, publishDate string, link string, releaseType models.ReleaseType) *models.Release {
	release := &models.Release{}
	release.Name = strings.TrimSpace(name)
	if release.Name == "" {
		return nil
	}

	for _, n := range artistNames {
		release.Artists = append(release.Artists, models.Artist{Name: strings.TrimSpace(n)})
	}

	date, err := time.Parse("January 2, 2006", publishDate)
	if err != nil {
		log.Printf("bandcamp: parse date %q: %v", publishDate, err)
		return nil
	}
	release.PublishDate = date
	release.Link = link
	release.ReleaseType = releaseType
	release.Source = models.SourceBandcamp

	return release
}

func GetAlbumOfTheDays(page int) []*models.Release {
	var albums []*models.Release

	c := colly.NewCollector()
	c.OnHTML(".list-article.aotd:not(.ft-latest-article)", func(e *colly.HTMLElement) {
		link := "https://daily.bandcamp.com" + e.ChildAttr("a", "href")
		date := e.ChildText(".article-info-text")
		date = strings.TrimSpace(strings.Split(date, "·")[1])

		title := normalizeTitle(e.ChildText(".title-wrapper"))
		match := bandcampTitleRe.FindStringSubmatch(title)

		if len(match) != 3 {
			log.Printf("bandcamp: unexpected title format %q", e.ChildText(".title-wrapper"))
			return
		}
		name := match[2]
		rawArtists := bandcampArtistSplitRe.Split(match[1], -1)
		var artists []string
		for _, a := range rawArtists {
			if a != "" {
				artists = append(artists, a)
			}
		}
		if len(artists) == 0 {
			return
		}

		r := makeBandcampRelease(name, artists, date, link, models.TypeAlbum)
		if r != nil {
			albums = append(albums, r)
		}
	})
	c.Visit("https://daily.bandcamp.com/album-of-the-day?page=" + strconv.Itoa(page))

	return albums
}
