package main

import (
	"fmt"
	"sync"
	"time"
	"os"
	"encoding/csv"
	"strings"
)

func main() {
	start := time.Now()
	releases := GetReleases()
	elapsed := time.Since(start)
	fmt.Println(len(releases))
	fmt.Printf("Total time: %s\n", elapsed)

	// Write releases to CSV
	file, err := os.Create("releases.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	file.Write(bomUtf8)

	// Write header
	writer.Write([]string{"Name", "Artists", "PublishDate", "Link", "Source", "ReleaseType", "Special"})

	for _, r := range releases {
		special := ""
		if r.Special != nil {
			special = *r.Special
		}
		writer.Write([]string{r.Name, strings.Join(r.Artists, ", "), r.PublishDate.String(), r.Link, r.Source, r.ReleaseType, special})
	}
}

func GetReleases() []Release {
	var wg sync.WaitGroup
	var mu sync.Mutex

	var releases []Release

	// bandcamp
	for i := range 100 {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			results := GetAlbumOfTheDays(page)
			mu.Lock()
			releases = append(releases, results...)
			mu.Unlock()
		}(i)
		// npr
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			results := GetNewMusicFriday(page, 10)
			mu.Lock()
			releases = append(releases, results...)
			mu.Unlock()
		}(i)
		// pitchfork
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			results := GetAlbums(page)
			mu.Lock()
			releases = append(releases, results...)
			mu.Unlock()
		}(i)
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			results := GetTracks(i)
			mu.Lock()
			releases = append(releases, results...)
			mu.Unlock()
		}(i)

		// stereogum
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			results := GetAlbumsOfTheWeek(page)
			mu.Lock()
			releases = append(releases, results...)
			mu.Unlock()
		}(i)

		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			results := GetAlbumsPremature(page)
			mu.Lock()
			releases = append(releases, results...)
			mu.Unlock()
		}(i)
		// xxl

		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			results := GetBangers(page, 12)
			mu.Lock()
			releases = append(releases, results...)
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return releases
}
