package models

import (
	"time"

	"gorm.io/gorm"
)

type ReleaseType string
type SpecialLabel string
type SourceName string

const (
	TypeAlbum ReleaseType = "album"
	TypeTrack ReleaseType = "track"
)

const (
	SpecialBestNewAlbum   SpecialLabel = "Best New Album"
	SpecialBestNewTrack   SpecialLabel = "Best New Track"
	SpecialBestNewReissue SpecialLabel = "Best New Reissue"
)

const (
	SourceBandcamp  SourceName = "Bandcamp"
	SourceNPR       SourceName = "NPR"
	SourcePitchfork SourceName = "Pitchfork"
	SourceStereogum SourceName = "Stereogum"
	SourceXXL       SourceName = "XXL"
)

type Artist struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex;not null"`
}

type Release struct {
	gorm.Model
	Name        string        `json:"name"`
	Artists     []Artist      `json:"artists" gorm:"many2many:release_artists;"`
	PublishDate time.Time     `json:"publish_date"`
	Link        string        `json:"link" gorm:"default:''"`
	Source      SourceName    `json:"source"`
	ReleaseType ReleaseType   `json:"release_type"`
	Special     *SpecialLabel `json:"special,omitempty"`
}
