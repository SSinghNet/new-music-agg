package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"log"

	"github.com/gocolly/colly"
)

func MakeNPRRelease(name string, artists []string, publishDate string, link string, releaseType string) Release {
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

	release.Source = "NPR"

	return release
}

// starts at page 0
func GetNewMusicFriday(page int, limit int) []Release {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
	)

	linkChan := make(chan string, limit+1)

	c.OnHTML("h2.title", func(e *colly.HTMLElement) {
		linkChan <- e.ChildAttr("a", "href")
	})

	c.OnScraped(func(r *colly.Response) {
		close(linkChan)
	})
	c.OnError(func(r *colly.Response, err error) {
		close(linkChan)
	})

	start := page*limit + 1
	url := fmt.Sprintf("https://www.npr.org/get/606254804/render/partial/next?start=%d&count=%d", start, limit)
	c.Visit(url)

	var albums []Release

	var wg sync.WaitGroup
	for link := range linkChan {
		wg.Add(1)
		go func(link string) {
			c := colly.NewCollector(
				colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
			)
			var date string
			c.OnHTML(".date", func(e *colly.HTMLElement) {
				date = e.Text
				date = strings.Replace(date, "Updated ", "", 1)
			})
			c.OnHTML(".edTag li", func(e *colly.HTMLElement) {
				exclude := []string{
					"Host: ",
					"Guest: ",
					"Producer: ",
					"Editor: ",
					"Executive Producer: ",
					"Vice President, Music and Visuals: ",
					"Vice President, Music & Visuals: ",
					"RIYL: ",
				}
				shouldExclude := false
				for _, substr := range exclude {
					if strings.Contains(e.Text, substr) {
						shouldExclude = true
						break
					}
				}
				if !shouldExclude {
					info := strings.Split(e.Text, ",")
					if len(info) == 2 {
						alb := MakeNPRRelease(info[1], strings.Split(info[0], "&"), date, link, "album")
						albums = append(albums, alb)
					}
				}
			})
			c.OnScraped(func(r *colly.Response) {
				wg.Done()
			})
			c.OnError(func(r *colly.Response, err error) {
				wg.Done()
			})
			c.Visit(link)
		}(link)
	}
	wg.Wait()
	
	return albums
}
