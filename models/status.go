package models

import (
	"time"
)

// Status represents the status of the last data refresh
type Status struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TotalCountries  int       `json:"total_countries"`
	LastRefreshedAt time.Time `json:"last_refreshed_at"`
}
