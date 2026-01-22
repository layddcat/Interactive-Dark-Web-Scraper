package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type DarkWebData struct {
	gorm.Model
	SourceName       string `json:"source_name"`
	SourceURL        string `json:"source_url"`
	RawContent       string `gorm:"type:text" json:"raw_content"`
	Title            string `json:"title"`
	PublishedAt      string `json:"published_at"`
	Category         string `json:"category"`
	CriticalityScore int    `json:"criticality_score"`
	IsReviewed       bool   `gorm:"default:false" json:"is_reviewed"`
}

type CategorySetting struct {
	gorm.Model
	Name            string `gorm:"unique"`
	BaseCriticality int
}

type IntelligenceData struct {
    ID               uint   `gorm:"primaryKey"`
    Title            string
    SourceName       string
    SourceURL        string
    RawContent       string
    Category         string
    CriticalityScore int
    PublishedAt      string
    IsActive         bool  
}