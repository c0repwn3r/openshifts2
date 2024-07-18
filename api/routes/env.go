package routes

import (
	"github.com/c0repwn3r/openshifts/api/apierr"
	"github.com/c0repwn3r/openshifts/api/config"
	"github.com/c0repwn3r/openshifts/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Env struct {
	Config config.Config
	DB     gorm.DB
}

func (e *Env) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("X-API-Token")
		if hdr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierr.ErrorResponse{Errors: []apierr.Error{
				{
					Code:    apierr.EUnauthorized,
					Message: "unauthorized",
					Path:    "",
				},
			}})
			return
		}
		token := &models.Token{}
		r := e.DB.Preload("User").Preload("User.Organization").Find(&token, "id = ?", hdr)
		if r.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierr.ErrorResponse{Errors: []apierr.Error{
				{
					Code:    apierr.EUnauthorized,
					Message: "unauthorized",
					Path:    "",
				},
			}})
			return
		}
		c.Set("user", token.User)
		c.Set("organization", token.User.Organization)
	}
}
