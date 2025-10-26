package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError represents a standardized error response
type APIError struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// NewAPIError creates a new APIError instance
func NewAPIError(err string, details interface{}) APIError {
	return APIError{
		Error:   err,
		Details: details,
	}
}

// HandleInternalServerError returns a 500 Internal Server Error response
func HandleInternalServerError(c *gin.Context, err error, msg string) {
	log.Printf("Internal Server Error: %s: %v", msg, err)
	c.JSON(http.StatusInternalServerError, NewAPIError("Internal server error", nil))
}

// HandleBadRequestError returns a 400 Bad Request error response
func HandleBadRequestError(c *gin.Context, details interface{}) {
	c.JSON(http.StatusBadRequest, NewAPIError("Validation failed", details))
}

// HandleNotFoundError returns a 404 Not Found error response
func HandleNotFoundError(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, NewAPIError(resource+" not found", nil))
}

// HandleServiceUnavailableError returns a 503 Service Unavailable error response
func HandleServiceUnavailableError(c *gin.Context, details string) {
	c.JSON(http.StatusServiceUnavailable, NewAPIError("External data source unavailable", details))
}
