package routes

import (
	"github.com/c0repwn3r/openshifts/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *Env) GetUserHoursRoute(c *gin.Context) {
	org := c.MustGet("user").(models.User)

	hours := models.GetUserHours(&e.DB, org.ID)

	c.JSON(http.StatusOK, gin.H{
		"hours": hours,
	})
}

type SetUserHoursRequest struct {
	Hours models.OrganizationHours `json:"hours"`
}

func (e *Env) SetUserHoursRoute(c *gin.Context) {
	var setHoursRequest SetUserHoursRequest
	if err := c.BindJSON(&setHoursRequest); err != nil {
		return
	}

	user := c.MustGet("user").(models.User)

	models.SetUserHours(&e.DB, user.ID, setHoursRequest.Hours)

	c.JSON(http.StatusOK, gin.H{
		"hours": setHoursRequest.Hours,
	})
}
