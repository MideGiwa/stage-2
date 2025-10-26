package models

import (
	"time"
)

// Country represents the structure of country data stored in the database
type Country struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"unique;not null" json:"name" binding:"required"`
	Capital         *string   `json:"capital"`
	Region          *string   `json:"region"`
	Population      uint64    `gorm:"not null" json:"population" binding:"required"`
	CurrencyCode    *string   `json:"currency_code" binding:"required"`
	ExchangeRate    *float64  `json:"exchange_rate"`
	EstimatedGDP    *float64  `json:"estimated_gdp"`
	FlagURL         *string   `json:"flag_url"`
	LastRefreshedAt time.Time `gorm:"autoUpdateTime" json:"last_refreshed_at"`
}
