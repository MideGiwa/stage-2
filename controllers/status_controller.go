package controllers

import (
	"net/http"

	"stage-2/services"
	"stage-2/utils"

	"github.com/gin-gonic/gin"
)

// StatusController handles HTTP requests related to the application status
type StatusController struct {
	statusService *services.StatusService
}

// NewStatusController creates a new StatusController
func NewStatusController(ss *services.StatusService) *StatusController {
	return &StatusController{statusService: ss}
}

// GetStatus handles the GET /status endpoint
func (ctrl *StatusController) GetStatus(c *gin.Context) {
	status, err := ctrl.statusService.GetStatus()
	if err != nil {
		utils.HandleInternalServerError(c, err, "failed to get status")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_countries":   status.TotalCountries,
		"last_refreshed_at": status.LastRefreshedAt.Format("2006-01-02T15:04:05Z"),
	})
}
