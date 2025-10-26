package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExternalAPIError is a custom error type for external API failures
type ExternalAPIError struct {
	Source string
	Err    error
}

func (e *ExternalAPIError) Error() string {
	return e.Source + " API failed: " + e.Err.Error()
}

// ExternalAPIErrorMiddleware catches ExternalAPIError and returns a 503 Service Unavailable response
func ExternalAPIErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				var extAPIErr *ExternalAPIError
				if errors.As(err.Err, &extAPIErr) {
					c.AbortWithStatusJSON(http.StatusServiceUnavailable, NewAPIError(
						"External data source unavailable",
						"Could not fetch data from "+extAPIErr.Source,
					))
					return
				}
			}
		}
	}
}
