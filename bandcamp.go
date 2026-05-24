package main

import (
	"log"
	"regexp"
	"strings"
	"time"
	"strconv"

	"github.com/gocolly/colly"
)

func makeBandcampRelease(name string, artists []string, publishDate string, link string, releaseType string) Release {
	release := Release{}
	release.Name = strings.TrimSpace(name)
	
	for i := range artists {
		artists[i] = strings.TrimSpace(artists[i])
	}
	release.Artists = artists

	date, err := time.Parse("January 2, 2006", publishDate)
	if err != nil {
		log.Fatal(err)
	}
	release.PublishDate = date

	release.Link = link
	release.ReleaseType = releaseType

	release.Source = "Bandcamp"

	return release
}


func GetAlbumOfTheDays(page int) ([]Release) {
	var albums []Release

	c := colly.NewCollector()
	c.OnHTML(".list-article.aotd:not(.ft-latest-article)", func(e *colly.HTMLElement) {
		link := "https://daily.bandcamp.com" + e.ChildAttr("a", "href")
		date := e.ChildText(".article-info-text")
		date = strings.TrimSpace(strings.Split(date, "·")[1])

		title := e.ChildText(".title-wrapper")
		re := regexp.MustCompile(`^(.*?)\s*[,，]\s*“([^”]+)”.*$`)
		match := re.FindStringSubmatch(title)

		var name string
		var artists []string
		if len(match) == 3 {
			name = match[2]
			splitArtist := regexp.MustCompile(`\s*(?:,|&)\s*`)
			artists = splitArtist.Split(match[1], -1)
		}

		albums = append(albums, makeBandcampRelease(name, artists, date, link, "album"))
	})
	c.Visit("https://daily.bandcamp.com/album-of-the-day?page="+strconv.Itoa(page))

	return albums
}
