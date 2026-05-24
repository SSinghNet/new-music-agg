package scraper

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/SSinghNet/new-music-agg/backend/internal/models"

	"github.com/gocolly/colly"
)

func makeStereogumRelease(name string, artistNames []string, link string, publishDate string, releaseType models.ReleaseType) *models.Release {
	release := &models.Release{}
	release.Name = strings.TrimSpace(name)
	if len(release.Name) == 0 {
		return nil
	}

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
		log.Printf("stereogum: parse date %q: %v", publishDate, err)
		return nil
	}
	release.PublishDate = date
	release.Link = link
	release.ReleaseType = releaseType
	release.Source = models.SourceStereogum

	return release
}

func scrapeStereogumCategory(baseURL string, labelPrefix string, releaseType models.ReleaseType, maxPages int) []*models.Release {
	var albums []*models.Release
	splitArtist := regexp.MustCompile(`\s*(?:,|&)\s*`)
	pagesFollowed := 0

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/115.0.0.0 Safari/537.36"),
	)

	c.OnHTML(`[class*="PostCard_stackedWrapper"]`, func(e *colly.HTMLElement) {
		link := "https://stereogum.com" + e.ChildAttr(`[class*="PostCard_titleLink"]`, "href")
		name := e.ChildText(`h3[class*="PostCard_title"] em`)
		if name == "" {
			name = e.ChildText(`h3[class*="PostCard_title"] i`)
		}

		h3Text := e.ChildText(`h3[class*="PostCard_title"]`)
		artist := strings.TrimSpace(strings.NewReplacer(
			labelPrefix, "",
			name, "",
		).Replace(h3Text))
		artist = strings.TrimRight(artist, " ,:")

		date := strings.TrimSpace(e.ChildText(`[class*="PostCard_dateWrapper"] span`))

		artists := splitArtist.Split(artist, -1)
		r := makeStereogumRelease(name, artists, link, date, releaseType)
		if r != nil {
			albums = append(albums, r)
		}
	})

	c.OnHTML(`[class*="Pagination_nextLink"]`, func(e *colly.HTMLElement) {
		if pagesFollowed < maxPages-1 {
			pagesFollowed++
			c.Visit("https://stereogum.com" + e.Attr("href"))
		}
	})

	c.Visit(baseURL)
	return albums
}

func GetAlbumsOfTheWeek(maxPages int) []*models.Release {
	return scrapeStereogumCategory(
		"https://stereogum.com/category/reviews/album-of-the-week/",
		"Album Of The Week:",
		models.TypeAlbum,
		maxPages,
	)
}

func GetAlbumsPremature(maxPages int) []*models.Release {
	return scrapeStereogumCategory(
		"https://stereogum.com/category/reviews/premature-evaluation/",
		"Premature Evaluation:",
		models.TypeAlbum,
		maxPages,
	)
}
