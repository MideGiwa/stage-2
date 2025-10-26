package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"stage-2/models"
	"stage-2/utils"

	"gorm.io/gorm"
)

// CountryService handles business logic related to countries
type CountryService struct {
	db         *gorm.DB
	httpClient *utils.HTTPClient
}

// NewCountryService creates a new CountryService
func NewCountryService(db *gorm.DB) *CountryService {
	return &CountryService{
		db:         db,
		httpClient: utils.NewHTTPClient(),
	}
}

// RefreshCountries fetches data from external APIs, processes it, and updates the database
func (s *CountryService) RefreshCountries() (int, time.Time, error) {
	var (
		countriesAPIResponse []struct {
			Name       string `json:"name"`
			Capital    string `json:"capital"`
			Region     string `json:"region"`
			Population uint64 `json:"population"`
			Flag       string `json:"flag"`
			Currencies []struct {
				Code string `json:"code"`
			} `json:"currencies"`
		}
		exchangeRatesAPIResponse struct {
			Rates map[string]float64 `json:"rates"`
		}
	)

	// Fetch countries data
	countriesURL := "https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies"
	if err := s.httpClient.Get(countriesURL, &countriesAPIResponse); err != nil {
		return 0, time.Time{}, &utils.ExternalAPIError{Source: "restcountries.com", Err: err}
	}
	log.Printf("Fetched %d countries from external API", len(countriesAPIResponse))

	// Fetch exchange rates
	exchangeRatesURL := "https://open.er-api.com/v6/latest/USD"
	if err := s.httpClient.Get(exchangeRatesURL, &exchangeRatesAPIResponse); err != nil {
		return 0, time.Time{}, &utils.ExternalAPIError{Source: "open.er-api.com", Err: err}
	}
	log.Printf("Fetched %d exchange rates from external API", len(exchangeRatesAPIResponse.Rates))

	now := time.Now().UTC()
	var processedCountries []models.Country

	// Process countries and calculate estimated GDP
	for _, apiCountry := range countriesAPIResponse {
		country := models.Country{
			Name:            apiCountry.Name,
			Capital:         &apiCountry.Capital,
			Region:          &apiCountry.Region,
			Population:      apiCountry.Population,
			FlagURL:         &apiCountry.Flag,
			LastRefreshedAt: now,
		}

		// Handle currency code
		if len(apiCountry.Currencies) > 0 {
			currencyCode := apiCountry.Currencies[0].Code
			country.CurrencyCode = &currencyCode

			// Match exchange rate
			if rate, ok := exchangeRatesAPIResponse.Rates[currencyCode]; ok {
				country.ExchangeRate = &rate
				// Compute estimated_gdp = population × random(1000–2000) ÷ exchange_rate
				randomMultiplier := float64(rand.Intn(1001) + 1000) // Random number between 1000 and 2000
				estimatedGDP := float64(country.Population) * randomMultiplier / rate
				country.EstimatedGDP = &estimatedGDP
			} else {
				log.Printf("Exchange rate for currency code %s not found for country %s. Setting exchange_rate and estimated_gdp to null.", currencyCode, country.Name)
				// Set to null as per requirements
				country.ExchangeRate = nil
				country.EstimatedGDP = nil
			}
		} else {
			log.Printf("No currency found for country %s. Setting currency_code, exchange_rate, and estimated_gdp to null/0.", country.Name)
			// Set to null/0 as per requirements
			country.CurrencyCode = nil
			country.ExchangeRate = nil
			estimatedGDP := 0.0 // Set to 0.0 for consistency if no currency
			country.EstimatedGDP = &estimatedGDP
		}

		processedCountries = append(processedCountries, country)
	}

	// Use a transaction for atomic updates/inserts
	tx := s.db.Begin()
	if tx.Error != nil {
		return 0, time.Time{}, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	for _, country := range processedCountries {
		var existingCountry models.Country
		// Case-insensitive comparison for name
		res := tx.Where("LOWER(name) = LOWER(?)", country.Name).First(&existingCountry)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				// Insert new record
				if err := tx.Create(&country).Error; err != nil {
					tx.Rollback()
					return 0, time.Time{}, fmt.Errorf("failed to insert country %s: %w", country.Name, err)
				}
			} else {
				// Other database error
				tx.Rollback()
				return 0, time.Time{}, fmt.Errorf("database error checking country %s: %w", country.Name, res.Error)
			}
		} else {
			// Update existing record
			country.ID = existingCountry.ID // Preserve ID for update
			if err := tx.Save(&country).Error; err != nil {
				tx.Rollback()
				return 0, time.Time{}, fmt.Errorf("failed to update country %s: %w", country.Name, err)
			}
		}
	}

	// Update global status
	var status models.Status
	// Always use ID 1 for the global status record
	res := tx.First(&status, 1)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return 0, time.Time{}, fmt.Errorf("failed to retrieve status: %w", res.Error)
	}

	status.ID = 1 // Ensure the ID is always 1 for this singleton status record
	status.TotalCountries = len(processedCountries)
	status.LastRefreshedAt = now

	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Create if not found
		if err := tx.Create(&status).Error; err != nil {
			tx.Rollback()
			return 0, time.Time{}, fmt.Errorf("failed to create status record: %w", err)
		}
	} else {
		// Update if found
		if err := tx.Save(&status).Error; err != nil {
			tx.Rollback()
			return 0, time.Time{}, fmt.Errorf("failed to update status record: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully refreshed %d countries in the database. Last refreshed at: %s", len(processedCountries), now.String())

	// Image Generation
	// Get top 5 countries by estimated GDP for image
	var allCountriesInDB []models.Country
	if err := s.db.Order("estimated_gdp DESC").Limit(5).Find(&allCountriesInDB).Error; err != nil {
		log.Printf("Warning: Failed to fetch top countries for image generation: %v", err)
		// Proceed without top countries if there's an error
	}

	if err := utils.GenerateSummaryImage(len(processedCountries), allCountriesInDB, now); err != nil {
		log.Printf("Warning: Failed to generate summary image: %v", err)
	}

	return len(processedCountries), now, nil
}

// GetCountries fetches all countries from the database with optional filters and sorting
func (s *CountryService) GetCountries(region, currency, sort string) ([]models.Country, error) {
	var countries []models.Country
	query := s.db.Model(&models.Country{})

	if region != "" {
		query = query.Where("LOWER(region) = LOWER(?)", region)
	}
	if currency != "" {
		query = query.Where("LOWER(currency_code) = LOWER(?)", currency)
	}

	// Default sort order
	orderBy := "name ASC"
	switch sort {
	case "gdp_desc":
		orderBy = "estimated_gdp DESC"
	case "gdp_asc":
		orderBy = "estimated_gdp ASC"
	case "name_desc":
		orderBy = "name DESC"
	case "name_asc":
		orderBy = "name ASC"
	case "population_desc":
		orderBy = "population DESC"
	case "population_asc":
		orderBy = "population ASC"
	}
	query = query.Order(orderBy)

	if err := query.Find(&countries).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}
	return countries, nil
}

// GetCountryByName fetches a single country by its name
func (s *CountryService) GetCountryByName(name string) (*models.Country, error) {
	var country models.Country
	// Case-insensitive search
	if err := s.db.Where("LOWER(name) = LOWER(?)", name).First(&country).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Country not found
		}
		return nil, fmt.Errorf("failed to fetch country by name %s: %w", name, err)
	}
	return &country, nil
}

// DeleteCountry deletes a country record by its name
func (s *CountryService) DeleteCountry(name string) error {
	result := s.db.Where("LOWER(name) = LOWER(?)", name).Delete(&models.Country{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete country %s: %w", name, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that no record was found to delete
	}
	return nil
}

// GetTopCountriesByGDP fetches the top N countries by estimated GDP
func (s *CountryService) GetTopCountriesByGDP(limit int) ([]models.Country, error) {
	var countries []models.Country
	if err := s.db.Order("estimated_gdp DESC").Limit(limit).Find(&countries).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch top countries by GDP: %w", err)
	}
	return countries, nil
}
