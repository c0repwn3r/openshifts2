package api

import (
	"github.com/c0repwn3r/openshifts/api/config"
	"github.com/c0repwn3r/openshifts/api/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Api struct {
	Config config.Config
	DB     gorm.DB
}

func (a Api) Run() error {
	log.Debug("building router")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	log.Debug("setting routes")

	env := &routes.Env{
		DB:     a.DB,
		Config: a.Config,
	}

	v1 := router.Group("/v1")
	{
		v1.POST("/user/register", env.RegisterRoute)

		authorized := v1.Group("/", env.AuthMiddleware())
		{
			org := authorized.Group("/org")
			{
				org.GET("/hours", env.GetHoursRoute)
				org.POST("/hours", env.SetHoursRoute)
			}
		}

	}

	log.WithFields(log.Fields{
		"listen": a.Config.Listen,
	}).Info("starting api")
	err := router.Run(a.Config.Listen)
	return err
}
