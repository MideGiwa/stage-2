package controllers

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"stage-2/models"
	"stage-2/services"
	"stage-2/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// CountryController handles HTTP requests related to countries
type CountryController struct {
	countryService *services.CountryService
	statusService  *services.StatusService
}

// NewCountryController creates a new CountryController
func NewCountryController(cs *services.CountryService, ss *services.StatusService) *CountryController {
	return &CountryController{countryService: cs, statusService: ss}
}

// RefreshCountries handles the POST /countries/refresh endpoint
func (ctrl *CountryController) RefreshCountries(c *gin.Context) {
	total, lastRefreshed, err := ctrl.countryService.RefreshCountries()
	if err != nil {
		// ExternalAPIError is handled by middleware
		utils.HandleInternalServerError(c, err, "failed to refresh countries")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Countries refreshed successfully",
		"total_countries":   total,
		"last_refreshed_at": lastRefreshed.Format("2006-01-02T15:04:05Z"),
	})
}

// GetCountries handles the GET /countries endpoint
func (ctrl *CountryController) GetCountries(c *gin.Context) {
	region := c.Query("region")
	currency := c.Query("currency")
	sort := c.Query("sort") // e.g., gdp_desc, name_asc, population_desc

	countries, err := ctrl.countryService.GetCountries(region, currency, sort)
	if err != nil {
		utils.HandleInternalServerError(c, err, "failed to get countries")
		return
	}

	c.JSON(http.StatusOK, countries)
}

// GetCountryByName handles the GET /countries/:name endpoint
func (ctrl *CountryController) GetCountryByName(c *gin.Context) {
	name := c.Param("name")

	country, err := ctrl.countryService.GetCountryByName(name)
	if err != nil {
		utils.HandleInternalServerError(c, err, "failed to get country by name")
		return
	}
	if country == nil {
		utils.HandleNotFoundError(c, "Country")
		return
	}

	c.JSON(http.StatusOK, country)
}

// DeleteCountry handles the DELETE /countries/:name endpoint
func (ctrl *CountryController) DeleteCountry(c *gin.Context) {
	name := c.Param("name")

	err := ctrl.countryService.DeleteCountry(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleNotFoundError(c, "Country")
			return
		}
		utils.HandleInternalServerError(c, err, "failed to delete country")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Country deleted successfully"})
}

// ServeSummaryImage handles the GET /countries/image endpoint
func (ctrl *CountryController) ServeSummaryImage(c *gin.Context) {
	imagePath := utils.GetSummaryImagePath()

	// Check if the file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, utils.NewAPIError("Summary image not found", nil))
		return
	}

	c.File(imagePath)
}

// validateCountry ensures the required fields for a country are present
func validateCountry(country *models.Country) map[string]string {
	validate := validator.New()
	err := validate.Struct(country)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			validationErrors[fieldName] = "is required"
		}
		return validationErrors
	}
	return nil
}
