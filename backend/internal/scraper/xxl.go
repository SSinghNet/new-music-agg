package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/SSinghNet/new-music-agg/backend/internal/models"

	"github.com/gocolly/colly"
)

func makeXXLRelease(name string, artist string, publishDate string, link string, releaseType models.ReleaseType) *models.Release {
	release := &models.Release{}
	release.Name = strings.TrimSpace(name)
	if release.Name == "" {
		return nil
	}

	splitArtist := regexp.MustCompile(`\s*(?:,|&|and)\s*`)
	for _, n := range splitArtist.Split(artist, -1) {
		trimmed := strings.TrimSpace(n)
		if trimmed != "" {
			release.Artists = append(release.Artists, models.Artist{Name: trimmed})
		}
	}
	if len(release.Artists) == 0 {
		return nil
	}

	date, err := time.Parse("2006-01-02 15:04:05 -0700", publishDate)
	if err != nil {
		log.Printf("xxl: parse date %q: %v", publishDate, err)
		return nil
	}
	release.PublishDate = date
	release.Link = "https:" + link
	release.ReleaseType = releaseType
	release.Source = models.SourceXXL

	return release
}

func GetBangers(page int, offset int) []*models.Release {
	url := fmt.Sprintf("https://www.xxlmag.com/rest/carbon/filter/main/tags/bangers/page/%d/offset/%d/", page, offset*page)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("xxl: fetch page %d: %v", page, err)
		return nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("xxl: read body page %d: %v", page, err)
		return nil
	}

	type xxlEntry struct {
		L3 struct {
			L4 struct {
				PublishDate string `json:"postDateGmt"`
				Title       string `json:"title"`
				Url         string `json:"url"`
			} `json:"mainData"`
		} `json:"data"`
	}
	type xxlResponse struct {
		L1 struct {
			L2 struct {
				Data map[string]xxlEntry `json:"dataDetails"`
			} `json:"carbonwidget/taxonomy-1"`
		} `json:"widgets"`
	}

	var xxl xxlResponse
	if err := json.Unmarshal(bodyBytes, &xxl); err != nil {
		log.Printf("xxl: unmarshal page %d: %v", page, err)
		return nil
	}

	var releases []*models.Release
	for _, v := range xxl.L1.L2.Data {
		title := v.L3.L4.Title
		publishDate := v.L3.L4.PublishDate
		u := v.L3.L4.Url
		if strings.Contains(title, "Songs") {
			releases = append(releases, getBangersTracks(u, publishDate)...)
		} else if strings.Contains(title, "Projects") {
			releases = append(releases, getBangersAlbums(u, publishDate)...)
		}
	}

	return releases
}

func getBangersTracks(url string, publishDate string) []*models.Release {
	var tracks []*models.Release

	c := colly.NewCollector()
	c.OnHTML("h3", func(h *colly.HTMLElement) {
		re := regexp.MustCompile(`^(.+?)(?:'s|')?\s*"([^"]+?)"(?:\s*[Ff]eaturing.*)?$`)
		match := re.FindStringSubmatch(h.Text)
		if len(match) == 3 {
			r := makeXXLRelease(match[2], match[1], publishDate, url, models.TypeTrack)
			if r != nil {
				tracks = append(tracks, r)
			}
		}
	})
	c.Visit(url)

	return tracks
}

func getBangersAlbums(url string, publishDate string) []*models.Release {
	var albums []*models.Release

	c := colly.NewCollector()
	c.OnHTML(".list-post-right header", func(h *colly.HTMLElement) {
		r := makeXXLRelease(h.ChildText("h2"), h.ChildText("small"), publishDate, url, models.TypeAlbum)
		if r != nil {
			albums = append(albums, r)
		}
	})
	c.Visit(url)

	return albums
}
