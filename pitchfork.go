package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func makePitchforkRelease(name string, artists []string, publishDate string, link string, special string, releaseType string) Release {
	release := Release{}

	name = strings.ReplaceAll(name, "“", "")
	name = strings.ReplaceAll(name, "”", "")
	name = strings.ReplaceAll(name, "’", "'")
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

	release.Link = "https://pitchfork.com" + link

	special = strings.TrimSpace(special)
	if special == "" {
		release.Special = nil
	} else {
		release.Special = &special
	}

	release.ReleaseType = releaseType

	release.Source = "Pitchfork"

	return release
}

func GetAlbums(page int) []Release {
	c := colly.NewCollector()
	var albums []Release
	c.OnHTML(".summary-item", func(e *colly.HTMLElement) {
		name := e.ChildText(".summary-item__hed")
		artists := strings.Split(e.ChildText(".summary-item__sub-hed"), "/")
		publishDate := e.ChildText(".summary-item__publish-date")
		link := e.ChildAttr(".summary-item__hed-link", "href")
		best := e.ChildText("[class^='SummaryItemReviewLabelWrapper-']")

		alb := makePitchforkRelease(name, artists, publishDate, link, best, "album")
		albums = append(albums, alb)
	})
	c.Visit("https://pitchfork.com/reviews/albums/?page="+strconv.Itoa(page))

	return albums
}

func GetTracks(page int) []Release {
	c := colly.NewCollector()
	var tracks []Release
	c.OnHTML(".summary-item", func(e *colly.HTMLElement) {
		name := e.ChildText(".summary-item__hed")
		artists := strings.Split(e.ChildText(".summary-item__sub-hed"), "/")
		publishDate := e.ChildText(".summary-item__publish-date")
		link := e.ChildAttr(".summary-item__hed-link", "href")
		best := e.ChildText("[class^='SummaryItemReviewLabelWrapper-']")

		track := makePitchforkRelease(name, artists, publishDate, link, best, "track")
		tracks = append(tracks, track)
	})
	c.Visit("https://pitchfork.com/reviews/tracks/?page="+strconv.Itoa(page))

	return tracks
}
