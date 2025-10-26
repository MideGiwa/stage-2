package services

import (
	"errors"
	"fmt"
	"time"

	"stage-2/models"

	"gorm.io/gorm"
)

// StatusService handles business logic related to the application status
type StatusService struct {
	db *gorm.DB
}

// NewStatusService creates a new StatusService
func NewStatusService(db *gorm.DB) *StatusService {
	return &StatusService{db: db}
}

// GetStatus retrieves the current application status from the database
func (s *StatusService) GetStatus() (*models.Status, error) {
	var status models.Status
	// Assume a single status record with ID 1
	if err := s.db.First(&status, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no status record exists, return a default empty status
			return &models.Status{
				TotalCountries:  0,
				LastRefreshedAt: time.Time{},
			}, nil
		}
		return nil, fmt.Errorf("failed to retrieve status: %w", err)
	}
	return &status, nil
}
