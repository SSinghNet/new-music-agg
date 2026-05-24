package main

import (
	"time"
	"fmt"
)

type Release struct {
	Name        string
	Artists     []string
	PublishDate time.Time
	Link        string
	Source 		string
	ReleaseType string // make enum
	Special     *string
}

func (r Release) String() string {
	// var special string
	// if r.Special != nil {
	// 	special = *r.Special
	// } else {
	// 	special = ""
	// }
	return fmt.Sprintf("%s: %s - %s, %s, %s, %s", r.ReleaseType, r.Name, r.Artists, r.PublishDate, r.Link, r.Source)
}