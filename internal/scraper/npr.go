package scraper

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/SSinghNet/new-music-agg/internal/models"

	"github.com/gocolly/colly"
)

func makeNPRRelease(name string, artistNames []string, publishDate string, link string, releaseType models.ReleaseType) *models.Release {
	release := &models.Release{}
	release.Name = strings.TrimSpace(name)
    publishDate = strings.TrimSpace(publishDate)

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
		log.Printf("npr: parse date %q: %v", publishDate, err)
		return nil
	}
	release.PublishDate = date
	release.Link = link
	release.ReleaseType = releaseType
	release.Source = models.SourceNPR

	return release
}

// GetNewMusicFriday starts at page 0.
func GetNewMusicFriday(page int, limit int) []*models.Release {
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

	var mu sync.Mutex
	var albums []*models.Release
	var wg sync.WaitGroup

	for link := range linkChan {
		wg.Add(1)
		go func(link string) {
			c2 := colly.NewCollector(
				colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
			)
			var date string
			c2.OnHTML(".date", func(e *colly.HTMLElement) {
				date = strings.Replace(e.Text, "Updated ", "", 1)
			})
			c2.OnHTML(".edTag li", func(e *colly.HTMLElement) {
				exclude := []string{
					"Host: ", "Guest: ", "Producer: ", "Editor: ",
					"Executive Producer: ",
					"Vice President, Music and Visuals: ",
					"Vice President, Music & Visuals: ",
					"RIYL: ",
                    "Recommended If You Like: ",
                    "Featured Songs: ",
                    "Featured Song: ",
				}
				for _, substr := range exclude {
					if strings.Contains(e.Text, substr) {
						return
					}
				}
				info := strings.Split(e.Text, ",")
				if len(info) == 2 {
					r := makeNPRRelease(info[1], strings.Split(info[0], "&"), date, link, models.TypeAlbum)
					if r != nil {
						mu.Lock()
						albums = append(albums, r)
						mu.Unlock()
					}
				}
			})
			c2.OnScraped(func(r *colly.Response) { wg.Done() })
			c2.OnError(func(r *colly.Response, err error) { wg.Done() })
			c2.Visit(link)
		}(link)
	}
	wg.Wait()

	return albums
}
