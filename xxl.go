package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func makeXXLRelease(name string, artist string, publishDate string, link string, releaseType string) Release {
	release := Release{}

	release.Name = strings.TrimSpace(name)

	var artists []string	
	splitArtist := regexp.MustCompile(`\s*(?:,|&|and)\s*`)
	artists = splitArtist.Split(artist, -1)

	for i := range artists {
		artists[i] = strings.TrimSpace(artists[i])
	}
	release.Artists = artists

	date, err := time.Parse("2006-01-02 15:04:05 -0700", publishDate)
	if err != nil {
		log.Fatal(err)
	}
	release.PublishDate = date

	release.Link = "https:" + link

	release.ReleaseType = releaseType

	return release
}

func GetBangers(page int, offset int) []Release {
	url := fmt.Sprintf("https://www.xxlmag.com/rest/carbon/filter/main/tags/bangers/page/%d/offset/%d/", page, offset*page)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type XXLResponse struct {
		L1 struct {
			L2 struct {
				Data map[string]struct {
					L3 struct {
						L4 struct {
							PublishDate string `json:"postDateGmt"`
							Title       string `json:"title"`
							Url         string `json:"url"`
						} `json:"mainData"`
					} `json:"data"`
				} `json:"dataDetails"`
			} `json:"carbonwidget/taxonomy-1"`
		} `json:"widgets"`
	}

	var xxl XXLResponse
	err = json.Unmarshal(bodyBytes, &xxl)
	if err != nil {
		log.Fatal(err)
	}


	var releases []Release
	for _, v := range xxl.L1.L2.Data {
		title := v.L3.L4.Title
		publishDate := v.L3.L4.PublishDate
		url := v.L3.L4.Url
		if strings.Contains(title, "Songs") {
			releases = append(releases, getBangersTracks(url, publishDate)...)
		} else if strings.Contains(title, "Projects") {
			releases = append(releases, getBangersAlbums(url, publishDate)...)
		}
	}

	return releases
}

func getBangersTracks(url string, publishDate string) []Release {
	var tracks []Release

	c := colly.NewCollector()
	c.OnHTML("h3", func(h *colly.HTMLElement) {
		text := h.Text
		re := regexp.MustCompile(`^(.+?)(?:'s|')?\s*"([^"]+?)"(?:\s*[Ff]eaturing.*)?$`)
		match := re.FindStringSubmatch(text)
		if len(match) == 3 {
			r := makeXXLRelease(match[2], match[1], publishDate, url, "track")
			tracks = append(tracks, r)
		} else {
			// fmt.Println("ERROR: ", text)
		}

	})

	c.Visit(url)
	
	return tracks
}

func getBangersAlbums(url string, publishDate string) []Release{
	var albums []Release

	c := colly.NewCollector()
	c.OnHTML(".list-post-right header", func(h *colly.HTMLElement) {
		name := h.ChildText("h2")
		artist := h.ChildText("small")
		// fmt.Println(name, "-", artist, publishDate)
		
		r := makeXXLRelease(name, artist, publishDate, url, "album")
		albums = append(albums, r)
	})
	c.Visit(url)

	return albums
}
