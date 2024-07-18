package routes

import (
	"github.com/c0repwn3r/openshifts/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *Env) GetHoursRoute(c *gin.Context) {
	org := c.MustGet("organization").(models.Organization)

	hours := models.GetOrgHours(&e.DB, org.ID)

	c.JSON(http.StatusOK, gin.H{
		"hours": hours,
	})
}

type SetHoursRequest struct {
	Hours models.OrganizationHours `json:"hours"`
}

func (e *Env) SetHoursRoute(c *gin.Context) {
	var setHoursRequest SetHoursRequest
	if err := c.BindJSON(&setHoursRequest); err != nil {
		return
	}

	org := c.MustGet("organization").(models.Organization)

	models.SetOrgHours(&e.DB, org.ID, setHoursRequest.Hours)

	c.JSON(http.StatusOK, gin.H{
		"hours": setHoursRequest.Hours,
	})
}
