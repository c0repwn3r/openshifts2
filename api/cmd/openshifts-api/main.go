package main

import (
	"flag"
	"fmt"
	"github.com/c0repwn3r/openshifts/api"
	"github.com/c0repwn3r/openshifts/api/config"
	"github.com/c0repwn3r/openshifts/api/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	configPath := flag.String("config", "", "Path to the configuration file")
	printUsage := flag.Bool("help", false, "Print usage")

	flag.Parse()

	if *printUsage {
		flag.Usage()
		os.Exit(0)
	}

	if *configPath == "" {
		fmt.Println("-config flag must be set")
		flag.Usage()
		os.Exit(1)
	}

	c, err := config.LoadConfig(*configPath)

	if err != nil {
		log.WithFields(log.Fields{
			"apierr": err,
		}).Error("configuration load failed")
		os.Exit(1)
	}

	log.WithFields(log.Fields{
		"c": c,
	}).Debug("config loaded")

	db, err := gorm.Open(postgres.Open(c.DbDSN), &gorm.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"apierr": err,
		}).Error("failed to connect to database")
		os.Exit(1)
	}

	log.Info("running migrations")
	err = db.AutoMigrate(&models.Availability{}, &models.Organization{}, &models.OrganizationWeekdayHours{}, &models.OrganizationOverrideHours{}, &models.OWHAvailability{}, &models.OOHAvailability{}, &models.User{}, &models.Token{}, &models.UserWeekdayHours{}, &models.UserOverrideHours{}, &models.UWHAvailability{}, &models.UOHAvailability{})
	if err != nil {
		log.WithFields(log.Fields{
			"apierr": err,
		}).Error("migrations failed")
		os.Exit(1)
	}

	a := api.Api{
		Config: *c,
		DB:     *db,
	}
	err = a.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"apierr": err,
		}).Error("api failed to run")
		os.Exit(1)
	}
}
