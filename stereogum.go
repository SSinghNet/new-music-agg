package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func MakeStereogumRelease(name string, artists []string, link string, publishDate string, releaseType string) Release {
	release := Release{}
	release.Name = strings.TrimSpace(name)

	for i := range artists {
		artists[i] = strings.TrimSpace(artists[i])
	}
	release.Artists = artists

	date, err := time.Parse("January 2, 2006", publishDate)
	if err != nil {
		log.Println(err)
	} else {
		release.PublishDate = date
	}

	release.Link = link

	release.ReleaseType = releaseType

	release.Source = "Stereogum"

	return release
}

func GetAlbumsOfTheWeek(page int) []Release {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var albums []Release

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
	)
	c.OnHTML(".article-card__title", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		aText := e.ChildText("a")
		name := e.ChildText("a em")
		if len(name) == 0 {
			name = e.ChildText("a i")
		}
		artist := strings.TrimSpace(strings.Replace(strings.Replace(aText, name, "", 1), "Album Of The Week:", "", 1))

		var artists []string
		splitArtist := regexp.MustCompile(`\s*(?:,|&)\s*`)
		artists = splitArtist.Split(artist, -1)
		wg.Add(1)

		go func(name string, artists []string, link string) {
			defer wg.Done()

			var date string

			c2 := colly.NewCollector(
				colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
			)
			c2.OnHTML(".date", func(f *colly.HTMLElement) {
				date = f.Text
			})
			c2.OnScraped(func(r *colly.Response) {
				mu.Lock()
				albums = append(albums, MakeStereogumRelease(name, artists, link, date, "album"))
				mu.Unlock()
			})
			c2.Visit(link)
		}(name, artists, link)

	})
	c.Visit("https://www.stereogum.com/category/reviews/album-of-the-week/page/" + strconv.Itoa(page))
	wg.Wait()

	return albums
}

//
// https://www.stereogum.com/category/reviews/premature-evaluation/page/2/

func GetAlbumsPremature(page int) []Release {
	var albums []Release

	var wg sync.WaitGroup
	var mu sync.Mutex

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
	)
	c.OnHTML(".article-card__title", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		aText := e.ChildText("a")
		name := e.ChildText("a em")
		if len(name) == 0 {
			name = e.ChildText("a i")
		}
		artist := strings.TrimSpace(strings.Replace(strings.Replace(aText, name, "", 1), "Premature Evaluation:", "", 1))

		var artists []string
		splitArtist := regexp.MustCompile(`\s*(?:,|&)\s*`)
		artists = splitArtist.Split(artist, -1)
		wg.Add(1)

		go func(name string, artists []string, link string) {
			defer wg.Done()
			var date string
			c2 := colly.NewCollector(
				colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
			)
			c2.OnHTML(".date", func(f *colly.HTMLElement) {
				date = f.Text
			})
			c2.OnScraped(func(r *colly.Response) {
				mu.Lock()
				albums = append(albums, MakeStereogumRelease(name, artists, link, date, "album"))
				mu.Unlock()
			})

			c2.Visit(link)

		}(name, artists, link)
	})
	c.Visit("https://www.stereogum.com/category/reviews/premature-evaluation/page/" + strconv.Itoa(page))
	wg.Wait()

	return albums
}
